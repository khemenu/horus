package provider

import (
	"golang.org/x/oauth2"
	"khepri.dev/horus"
)

type provider struct {
	id horus.Verifier
}

func (p provider) Id() horus.Verifier {
	return p.id
}

type OauthProviderConfig struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type oauthProvider struct {
	provider
	conf oauth2.Config
}

func (p *oauthProvider) Config() *oauth2.Config {
	return &p.conf
}
