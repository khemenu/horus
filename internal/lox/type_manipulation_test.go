package lox_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"khepri.dev/horus/internal/lox"
)

func TestFromPtrOrF(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const fallbackStr = "fallback"
	callbackStr := func() string {
		return fallbackStr
	}
	str := "foo"
	ptrStr := &str

	const fallbackInt = -1
	callbackInt := func() int {
		return fallbackInt
	}
	i := 9
	ptrInt := &i

	is.Equal(str, lox.FromPtrOrF(ptrStr, callbackStr))
	is.Equal(fallbackStr, lox.FromPtrOrF(nil, callbackStr))
	is.Equal(i, lox.FromPtrOrF(ptrInt, callbackInt))
	is.Equal(fallbackInt, lox.FromPtrOrF(nil, callbackInt))
}
