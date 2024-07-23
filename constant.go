package horus

const (
	TokenTypePassword = "password"
	TokenTypeRefresh  = "refresh"
	TokenTypeAccess   = "access"
	TokenTypeOtp      = "otp"
)

const (
	TokenKeyName = "horus_token"
)

const (
	// Name of the silo that holds the permission for `Conf`.
	ConfSiloName = "_config"
)

const (
	InvitationTypeInternal = "internal"
)

const (
	// Me is an alias that indicates a resource the user owns.
	Me = "_me"
)
