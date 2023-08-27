package provider

import (
	"context"
	"net/http"
	"net/http/httptest"

	"golang.org/x/oauth2"
	"khepri.dev/horus"
)

type fakeOauthProvider struct {
	oauthProvider
}

func FakeOauth2() (horus.OauthProvider, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"access_token": "90d64460d14870c08c81352a05dedd3465940a7c",
			"token_type": "Bearer",
			"expires_in": 3600
		}`))
	}))

	return &fakeOauthProvider{
		oauthProvider: oauthProvider{
			provider: provider{id: horus.VerifierFakeOauth2},
			conf: oauth2.Config{
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://example.com/authorize",
					TokenURL: server.URL,
				},
			},
		},
	}, server
}

func (p *fakeOauthProvider) Identity(ctx context.Context, token *oauth2.Token) (horus.IdentityInit, error) {
	return horus.IdentityInit{
		Value:      "ra@example,com",
		Kind:       horus.IdentityMail,
		VerifiedBy: p.id,
	}, nil
}
