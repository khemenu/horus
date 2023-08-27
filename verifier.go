package horus

type Verifier string

const (
	Unverified           Verifier = "unverified"
	VerifierFakeOauth2   Verifier = "fake-oauth2"
	VerifierGoogleOauth2 Verifier = "google-oauth2"
)
