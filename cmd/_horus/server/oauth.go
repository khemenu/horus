package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"khepri.dev/horus"
)

type oauthState struct {
	Nonce      string `json:"nonce"`
	ProviderId string `json:"provider_id"`
	RedirectTo string `json:"redirect_to"`
}

func (s *restServer) OauthRedirect(w http.ResponseWriter, r *http.Request) {
	state := oauthState{
		ProviderId: r.URL.Query().Get("provider_id"),
		RedirectTo: r.URL.Query().Get("redirect_to"),
	}

	provider, ok := s.oauth_providers[horus.Verifier(state.ProviderId)]
	if !ok {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	opaque, err := horus.DefaultOpaqueTokenGenerator.New()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	state.Nonce = opaque
	state_json, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}

	state_str := url.QueryEscape(string(state_json))

	// Set OAuth state
	http.SetCookie(w, &http.Cookie{
		Name:  horus.CookieNameOauthState,
		Value: state_str,

		Path:    filepath.Join(s.conf.AppPrefix, "oauth/callback"),
		Domain:  s.conf.AppDomain,
		Expires: time.Now().Add(5 * time.Minute),

		Secure:   s.conf.Debug.IsSecure(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	next := provider.Config().AuthCodeURL(state_str, oauth2.AccessTypeOnline)
	http.Redirect(w, r, next, http.StatusTemporaryRedirect)
}

func (s *restServer) OauthCallback(w http.ResponseWriter, r *http.Request) {
	actual := ""
	if v, err := url.QueryUnescape(r.FormValue("state")); err == nil {
		actual = v
	}

	expected := ""
	if cookie, err := r.Cookie(horus.CookieNameOauthState); err == nil {
		if v, err := url.QueryUnescape(cookie.Value); err == nil {
			expected = v
		}
	}

	if actual == "" || expected == "" || actual != expected {
		http.Error(w, "OAuth state does not match", http.StatusUnprocessableEntity)
		return
	}

	state := oauthState{}
	if err := json.Unmarshal([]byte(expected), &state); err != nil {
		http.Error(w, "invalid OAuth state", http.StatusBadRequest)
		return
	}

	provider, ok := s.oauth_providers[horus.Verifier(state.ProviderId)]
	if !ok {
		// Process restarted or deployment being updated after removing the provider from the config.
		http.Error(w, "unknown OAuth provider", http.StatusInternalServerError)
		return
	}

	// Invalidate OAuth state
	http.SetCookie(w, &http.Cookie{
		Name:   horus.CookieNameOauthState,
		Value:  "",
		MaxAge: -1,
	})

	code := r.FormValue("code")
	token, err := provider.Config().Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "failed to exchange OAuth token", http.StatusInternalServerError)
		return
	}

	identity_init, err := provider.Identity(r.Context(), token)
	if err != nil {
		http.Error(w, "failed to resolve identity", http.StatusInternalServerError)
		return
	}

	var user_id horus.UserId
	if identity_registered, err := s.Identities().GetByValue(r.Context(), identity_init.Value); err != nil {
		if !errors.Is(err, horus.ErrNotExist) {
			http.Error(w, "failed to get identity from store", http.StatusInternalServerError)
			return
		}

		identity, err := s.Identities().New(r.Context(), &identity_init)
		if err != nil {
			http.Error(w, "failed to sign up", http.StatusInternalServerError)
			return
		}

		user, err := s.Users().GetById(r.Context(), identity.OwnerId)
		if err != nil {
			http.Error(w, "failed to get details of created identity", http.StatusInternalServerError)
			return
		}

		user_id = user.Id
	} else {
		user_id = identity_registered.OwnerId
	}

	if _, ok := s.resetRefreshToken(w, r, user_id); !ok {
		return
	}
	if _, ok := s.resetAccessToken(w, r, user_id); !ok {
		return
	}

	if state.RedirectTo == "" {
		state.RedirectTo = "/"
	}
	http.Redirect(w, r, state.RedirectTo, http.StatusTemporaryRedirect)
}
