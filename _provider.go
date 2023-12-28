package horus

import (
	"context"

	"golang.org/x/oauth2"
)

type Verifier string

const (
	Unverified           Verifier = "unverified"
	VerifierFakeOauth2   Verifier = "fake-oauth2"
	VerifierGoogleOauth2 Verifier = "google-oauth2"
)

type Provider interface {
	Id() Verifier
}

type OauthProvider interface {
	Provider
	Config() *oauth2.Config
	Identity(ctx context.Context, token *oauth2.Token) (IdentityInit, error)
}
