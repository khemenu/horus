package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"khepri.dev/horus"
)

type googleOauthProvider struct {
	oauthProvider
}

func GoogleOauth2(conf OauthProviderConfig) horus.OauthProvider {
	return &googleOauthProvider{
		oauthProvider: oauthProvider{
			provider: provider{id: horus.VerifierGoogleOauth2},
			conf: oauth2.Config{
				ClientID:     conf.ClientId,
				ClientSecret: conf.ClientSecret,
				Endpoint:     google.Endpoint,
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
				},
			},
		},
	}
}

func (p *googleOauthProvider) Identity(ctx context.Context, token *oauth2.Token) (horus.IdentityInit, error) {
	raw := token.Extra("id_token")
	entires := strings.SplitN(raw.(string), ".", 3)
	if len(entires) != 3 {
		return horus.IdentityInit{}, errors.New("expected a JWT token")
	}

	entry := entires[1]
	payload, err := base64.RawStdEncoding.DecodeString(entry)
	if err != nil {
		return horus.IdentityInit{}, fmt.Errorf("decode payload: %w", err)
	}

	id_token := struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}{}
	if err := json.Unmarshal(payload, &id_token); err != nil {
		return horus.IdentityInit{}, fmt.Errorf("expected payload to be a valid JSON: %w", err)
	}
	if !id_token.EmailVerified {
		return horus.IdentityInit{}, errors.New("email not verified")
	}

	return horus.IdentityInit{
		Value: id_token.Email,

		Kind:       horus.IdentityEmail,
		VerifiedBy: p.id,
	}, nil
}
