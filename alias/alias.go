package alias

import (
	"errors"
	"math/rand"
	"regexp"
)

const Charset = "abcdefghijklmnopqrstuvwxyz0123456789"

var (
	BeginAlphabetEndNonSymbol = regexp.MustCompile(`^[a-z](?:.*[a-z0-9])?$`)
	ContainsNotAllowedChars   = regexp.MustCompile(`[^a-z0-9-_]`)
	ConsecutiveSymbols        = regexp.MustCompile(`--|__`)
)

func New() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = Charset[rand.Intn(len(Charset))]
	}
	b[0] = Charset[rand.Intn('z'-'a')]
	return string(b)
}

func ValidateE(s string) error {
	switch {
	case !BeginAlphabetEndNonSymbol.MatchString(s):
		return errors.New("it must begin with a lowercase alphabet and end with a lowercase alphabet or a number")
	case ContainsNotAllowedChars.MatchString(s):
		return errors.New("only lowercase alphabet, '-', and '_' characters are allowed")
	case ConsecutiveSymbols.MatchString(s):
		return errors.New("consecutive '-' or '_' characters are not allowed")
	default:
		return nil
	}
}

func Validate(s string) bool {
	return ValidateE(s) == nil
}
