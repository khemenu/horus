package horus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var (
	CookieNameRefreshToken = "horus_refresh_token"
	CookieNameAccessToken  = "horus_access_token"
	CookieNameOauthState   = "horus_oauth_state"
)

type Config struct {
	AppDomain string
	AppPrefix string
}

func (c *Config) Normalize() error {
	errs := []error{}

	if _, err := url.Parse(c.AppDomain); err != nil {
		errs = append(errs, fmt.Errorf("domain must be a valid URL"))
	}
	if c.AppPrefix == "" {
		c.AppPrefix = "/auth"
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

type Horus interface {
	Stores
	Services

	Verify(next http.HandlerFunc) http.HandlerFunc
	// VerifyF(next http.HandlerFunc) http.Handler

	Config() *Config
}

type Stores interface {
	Users() UserStore
	Tokens() TokenStore
	Identities() IdentityStore

	Orgs() OrgStore
	Teams() TeamStore
	Members() MemberStore
	Memberships() MembershipStore
}

type Services interface {
	Auth() AuthService
}

type AuthService interface {
	SignUp(ctx context.Context, init IdentityInit) (*User, error)
}
