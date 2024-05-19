package fx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"khepri.dev/horus/internal/fx"
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

	is.Equal(str, fx.FromPtrOrF(ptrStr, callbackStr))
	is.Equal(fallbackStr, fx.FromPtrOrF(nil, callbackStr))
	is.Equal(i, fx.FromPtrOrF(ptrInt, callbackInt))
	is.Equal(fallbackInt, fx.FromPtrOrF(nil, callbackInt))
}
