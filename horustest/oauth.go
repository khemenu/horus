package horustest

import (
	"testing"

	"khepri.dev/horus"
	"khepri.dev/horus/provider"
)

func WithFakeOauth2(f func(provider horus.OauthProvider)) func(t *testing.T) {
	return func(t *testing.T) {
		provider, server := provider.FakeOauth2()
		defer server.Close()

		f(provider)
	}
}
