package conf

type ConfSignInLockout struct {
	Enabled bool
	Count   uint

	// Lock the user for this value of minutes.
	LockedPeriod uint `json:"locked_period"`
}

func (ConfSignInLockout) Id() string {
	return "signIn.lockout"
}
