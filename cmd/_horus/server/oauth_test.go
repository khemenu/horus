package server_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/horus/server"
	"khepri.dev/horus/provider"
)

func WithFakeOauth2(f func(provider horus.OauthProvider)) func(t *testing.T) {
	return func(t *testing.T) {
		provider, server := provider.FakeOauth2()
		defer server.Close()

		f(provider)
	}
}

func TestOauth(t *testing.T) {
	t.Run("redirect without provider id", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		require.HTTPStatusCode(handler.ServeHTTP, http.MethodGet, "/auth/oauth/redirect", nil, http.StatusUnprocessableEntity)
	}))

	t.Run("redirect with unknown provider id", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		require.HTTPStatusCode(handler.ServeHTTP, http.MethodGet, "/auth/oauth/redirect?provider_id=not-exist", nil, http.StatusUnprocessableEntity)
	}))

	t.Run("redirect with known provider id", WithFakeOauth2(func(provider horus.OauthProvider) {
		conf := server.RestServerConfig{
			Providers: []horus.Provider{provider},
		}

		WithHorusHttpHandler(&conf, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf(`/auth/oauth/redirect?provider_id=%s`, provider.Id()), nil)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			res := w.Result()
			require.Equal(http.StatusTemporaryRedirect, res.StatusCode)

			url, err := res.Location()
			require.NoError(err)

			url.RawQuery = ""
			require.Equal(provider.Config().Endpoint.AuthURL, url.String())
		})(t)
	}))

	t.Run("callback without state", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		require.HTTPStatusCode(handler.ServeHTTP, http.MethodGet, "/auth/oauth/callback", nil, http.StatusUnprocessableEntity)
	}))

	t.Run("callback with unknown provider id", WithFakeOauth2(func(provider horus.OauthProvider) {
		conf := server.RestServerConfig{
			Providers: []horus.Provider{provider},
		}

		WithHorusHttpHandler(&conf, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
			state := url.QueryEscape(`{"provider_id":"not_exist"}`)

			r := httptest.NewRequest(http.MethodGet, "/auth/oauth/callback", nil)
			r.AddCookie(&http.Cookie{Name: horus.CookieNameOauthState, Value: state})
			r.Form = url.Values{"state": []string{state}}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			res := w.Result()
			require.Equal(http.StatusInternalServerError, res.StatusCode)
		})(t)
	}))

	t.Run("callback with valid state", WithFakeOauth2(func(provider horus.OauthProvider) {
		conf := server.RestServerConfig{
			Providers: []horus.Provider{provider},
		}

		WithHorusHttpHandler(&conf, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
			state := url.QueryEscape(fmt.Sprintf(`{"provider_id":"%s"}`, provider.Id()))

			r := httptest.NewRequest(http.MethodGet, "/auth/oauth/callback", nil)
			r.AddCookie(&http.Cookie{Name: horus.CookieNameOauthState, Value: state})
			r.Form = url.Values{"state": []string{state}}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			res := w.Result()
			require.Equal(http.StatusTemporaryRedirect, res.StatusCode)

			cookie_names := []string{}
			for _, cookie := range res.Cookies() {
				cookie_names = append(cookie_names, cookie.Name)
			}
			require.Contains(cookie_names, horus.CookieNameRefreshToken)
			require.Contains(cookie_names, horus.CookieNameAccessToken)
		})(t)
	}))
}
