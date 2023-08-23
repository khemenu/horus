package server_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"khepri.dev/horus"
	"khepri.dev/horus/horustest"
	"khepri.dev/horus/pb"
	"khepri.dev/horus/server"
)

func WithHorusHttpHandler(conf *server.RestServerConfig, f func(require *require.Assertions, h horus.Horus, handler http.Handler)) func(t *testing.T) {
	if conf == nil {
		conf = &server.RestServerConfig{}
	}
	return horustest.WithHorus(conf.Config, func(require *require.Assertions, h horus.Horus) {
		s, err := server.NewRestServer(h, conf)
		require.NoError(err)

		f(require, h, s)
	})
}

func TestNotFound(t *testing.T) {
	WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		require.HTTPStatusCode(handler.ServeHTTP, http.MethodGet, "not_exist", nil, http.StatusNotFound)
		require.HTTPStatusCode(handler.ServeHTTP, http.MethodGet, "/auth/not_exist", nil, http.StatusNotFound)
	})(t)
}

func TestAuth(t *testing.T) {
	t.Run("sign out without tokens", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		require.HTTPStatusCode(handler.ServeHTTP, http.MethodGet, "/auth/signout", nil, http.StatusOK)
	}))

	t.Run("sign out with invalid tokens", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		r := httptest.NewRequest(http.MethodGet, "/auth/signout", nil)
		r.AddCookie(&http.Cookie{Name: horus.CookieNameRefreshToken, Value: "foo"})
		r.AddCookie(&http.Cookie{Name: horus.CookieNameAccessToken, Value: "bar"})
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		res := w.Result()
		require.Equal(http.StatusOK, res.StatusCode)

		cookies := map[string]*http.Cookie{}
		for _, cookie := range res.Cookies() {
			cookies[cookie.Name] = cookie
		}
		require.Contains(cookies, horus.CookieNameRefreshToken)
		require.Contains(cookies, horus.CookieNameAccessToken)
		require.Empty(cookies[horus.CookieNameRefreshToken].Value)
		require.Empty(cookies[horus.CookieNameAccessToken].Value)
		require.LessOrEqual(cookies[horus.CookieNameRefreshToken].MaxAge, 0)
		require.LessOrEqual(cookies[horus.CookieNameAccessToken].MaxAge, 0)
	}))

	t.Run("refresh without refresh token", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		require.HTTPStatusCode(handler.ServeHTTP, http.MethodGet, "/auth/refresh", nil, http.StatusUnprocessableEntity)
	}))

	t.Run("refresh with invalid refresh token", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		r := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
		r.AddCookie(&http.Cookie{Name: horus.CookieNameRefreshToken, Value: "foo"})
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		res := w.Result()
		require.Equal(http.StatusUnauthorized, res.StatusCode)
	}))

	t.Run("refresh with valid refresh token", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		ctx := context.Background()

		user, err := h.Users().New(ctx)
		require.NoError(err)

		refresh_token, err := h.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  user.Id,
			Type:     horus.RefreshToken,
			Duration: time.Hour,
		})
		require.NoError(err)

		r := httptest.NewRequest(http.MethodGet, "/auth/refresh", nil)
		r.AddCookie(&http.Cookie{Name: horus.CookieNameRefreshToken, Value: refresh_token.Value})
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		res := w.Result()
		require.Equal(http.StatusOK, res.StatusCode)
		require.True(slices.ContainsFunc(res.Cookies(), func(c *http.Cookie) bool {
			return c.Name == horus.CookieNameAccessToken
		}))

		body, err := io.ReadAll(res.Body)
		require.NoError(err)

		msg := pb.StatusRes{}
		err = protojson.Unmarshal(body, &msg)
		require.NoError(err)
	}))

	t.Run("status without access token", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		require.HTTPStatusCode(handler.ServeHTTP, http.MethodGet, "/auth/status", nil, http.StatusUnprocessableEntity)
	}))

	t.Run("status with invalid access token", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		r := httptest.NewRequest(http.MethodGet, "/auth/status", nil)
		r.AddCookie(&http.Cookie{Name: horus.CookieNameAccessToken, Value: "foo"})
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		res := w.Result()
		require.Equal(http.StatusUnprocessableEntity, res.StatusCode)
	}))

	t.Run("status with valid access token", WithHorusHttpHandler(nil, func(require *require.Assertions, h horus.Horus, handler http.Handler) {
		ctx := context.Background()

		user, err := h.Users().New(ctx)
		require.NoError(err)

		access_token, err := h.Tokens().Issue(ctx, horus.TokenInit{
			OwnerId:  user.Id,
			Type:     horus.AccessToken,
			Duration: time.Hour,
		})
		require.NoError(err)

		r := httptest.NewRequest(http.MethodGet, "/auth/status", nil)
		r.AddCookie(&http.Cookie{Name: horus.CookieNameAccessToken, Value: access_token.Value})
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		res := w.Result()
		require.Equal(http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		require.NoError(err)

		msg := pb.StatusRes{}
		err = protojson.Unmarshal(body, &msg)
		require.NoError(err)
	}))
}
