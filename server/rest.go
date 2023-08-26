package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"khepri.dev/horus"
	"khepri.dev/horus/frame"
	"khepri.dev/horus/log"
	"khepri.dev/horus/pb"
)

type RestServerDebugConfig struct {
	Enabled   bool
	Unsecured bool
}

func (c *RestServerDebugConfig) IsSecure() bool {
	if !c.Enabled {
		return true
	}

	return !c.Unsecured
}

type RestServerConfig struct {
	*horus.Config

	Providers []horus.Provider

	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration

	Debug RestServerDebugConfig
}

func (c *RestServerConfig) Normalize() error {
	errs := []error{}

	if c.Config == nil {
		c.Config = &horus.Config{}
	}
	if err := c.Config.Normalize(); err != nil {
		errs = append(errs, err)
	}

	if c.AccessTokenExpiry == 0 {
		c.AccessTokenExpiry = 6 * time.Hour
	}
	if c.RefreshTokenExpiry == 0 {
		c.RefreshTokenExpiry = 6 * 30 * 24 * time.Hour // About 6 months.
	}

	if len(errs) > 0 {
		return fmt.Errorf("invalid config: %w", errors.Join(errs...))
	}

	return nil
}

type restServer struct {
	horus.Horus
	oauth_providers map[horus.Verifier]horus.OauthProvider

	conf *RestServerConfig
}

func NewRestServer(h horus.Horus, conf *RestServerConfig) (http.Handler, error) {
	if conf == nil {
		conf = &RestServerConfig{}
	}
	if conf.Config == nil {
		conf.Config = h.Config()
	}
	if err := conf.Normalize(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	oauth_redirect_url, err := url.Parse(conf.AppDomain)
	if err != nil {
		panic("invalid URL")
	}

	oauth_redirect_url.Path = filepath.Join(conf.AppPrefix, "oauth/callback")

	oauth_providers := map[horus.Verifier]horus.OauthProvider{}
	for _, provider := range conf.Providers {
		switch p := provider.(type) {
		case horus.OauthProvider:
			p.Config().RedirectURL = oauth_redirect_url.String()
			oauth_providers[p.Id()] = p
		}
	}

	return &restServer{
		Horus:           h,
		oauth_providers: oauth_providers,

		conf: conf,
	}, nil
}

func (s *restServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, s.conf.AppPrefix) {
		http.NotFound(w, r)
		return
	}

	http.StripPrefix(s.conf.AppPrefix, http.HandlerFunc(s.handleRoot)).ServeHTTP(w, r)
}

func (s *restServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/whoami":
		if s.conf.Debug.Enabled {
			s.Verify(s.WhoAmI).ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}

	case "/signout":
		s.SignOut(w, r)

	case "/refresh":
		s.Refresh(w, r)

	case "/status":
		s.Status(w, r)

	case "/oauth/redirect":
		s.OauthRedirect(w, r)

	case "/oauth/callback":
		s.OauthCallback(w, r)

	default:
		http.NotFound(w, r)
		return
	}
}

func (s *restServer) WhoAmI(w http.ResponseWriter, r *http.Request) {
	frame := frame.MustFromCtx(r.Context())

	user, err := frame.User(r.Context())
	if err != nil {
		http.Error(w, "failed to get user details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"user_id":"%s", "alias":"%s", "created_at":"%s"}`, user.Id, user.Alias, user.CreatedAt.Format(time.RFC3339))
}

func (s *restServer) SignOut(w http.ResponseWriter, r *http.Request) {
	l := log.FromCtx(r.Context())
	if cookie, err := r.Cookie(horus.CookieNameRefreshToken); err == nil {
		if err := s.Tokens().Revoke(r.Context(), cookie.Value); err != nil {
			l.Warn("failed to revoke a token", "err", err)
		}
	}
	if cookie, err := r.Cookie(horus.CookieNameAccessToken); err == nil {
		if err := s.Tokens().Revoke(r.Context(), cookie.Value); err != nil {
			l.Warn("failed to revoke a token", "err", err)
		}
	}

	http.SetCookie(w, &http.Cookie{Name: horus.CookieNameRefreshToken, Value: "", MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: horus.CookieNameAccessToken, Value: "", Path: "/", MaxAge: -1})

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"msg":"live long and prosper"}`)
}

func (s *restServer) resetRefreshToken(w http.ResponseWriter, r *http.Request, user_id horus.UserId) (*horus.Token, bool) {
	refresh_token, err := s.Tokens().Issue(r.Context(), horus.TokenInit{
		OwnerId:  user_id,
		Type:     horus.RefreshToken,
		Duration: s.conf.RefreshTokenExpiry,
	})
	if err != nil {
		http.Error(w, "issue refresh token", http.StatusInternalServerError)
		return nil, false
	}

	http.SetCookie(w, &http.Cookie{
		Name:  horus.CookieNameRefreshToken,
		Value: refresh_token.Value,

		Path:    filepath.Join(s.conf.AppPrefix, "refresh"),
		Domain:  s.conf.AppDomain,
		Expires: refresh_token.ExpiredAt,

		Secure:   s.conf.Debug.IsSecure(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return refresh_token, true
}

func (s *restServer) resetAccessToken(w http.ResponseWriter, r *http.Request, user_id horus.UserId) (*horus.Token, bool) {
	access_token, err := s.Tokens().Issue(r.Context(), horus.TokenInit{
		OwnerId:  user_id,
		Type:     horus.AccessToken,
		Duration: s.conf.AccessTokenExpiry,
	})
	if err != nil {
		http.Error(w, "failed to issue access token", http.StatusInternalServerError)
		return nil, false
	}

	http.SetCookie(w, &http.Cookie{
		Name:  horus.CookieNameAccessToken,
		Value: access_token.Value,

		Path:    "/",
		Domain:  s.conf.AppDomain,
		Expires: access_token.ExpiredAt,

		Secure:   s.conf.Debug.IsSecure(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return access_token, true
}

func (s *restServer) status(w http.ResponseWriter, access_token *horus.Token) {
	rst, err := protojson.Marshal(&pb.StatusRes{
		SessionExpiredAt: access_token.ExpiredAt.Format(time.RFC3339),
	})
	if err != nil {
		http.Error(w, "failed to marshal the result", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(rst)
}

func (s *restServer) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(horus.CookieNameRefreshToken)
	if err != nil {
		http.Error(w, "no refresh token", http.StatusUnprocessableEntity)
		return
	}

	refresh_token, err := s.Tokens().GetByValue(r.Context(), cookie.Value, horus.RefreshToken)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		http.Error(w, "failed to get token details", http.StatusInternalServerError)
		return
	}

	access_token, ok := s.resetAccessToken(w, r, refresh_token.OwnerId)
	if !ok {
		return
	}

	s.status(w, access_token)
}

func (s *restServer) Status(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(horus.CookieNameAccessToken)
	if err != nil {
		http.Error(w, "no access token", http.StatusUnprocessableEntity)
		return
	}

	access_token, err := s.Tokens().GetByValue(r.Context(), cookie.Value, horus.AccessToken)
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			// Try refresh
			s.Refresh(w, r)
			return
		}

		http.Error(w, "failed to get token details", http.StatusInternalServerError)
		return
	}

	s.status(w, access_token)
}
