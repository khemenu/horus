package fx_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus/internal/fx"
)

func TestCollectErr(t *testing.T) {
	require := require.New(t)

	err_foo := errors.New("foo")
	err_bar := errors.New("bar")

	a := func() (int, error) { return 42, err_foo }
	b := func() (float32, error) { return 3.14, err_bar }

	errs := []error{}

	answer := fx.CollectErr(a()).To(&errs)
	pi := fx.CollectErr(b()).To(&errs)
	require.Equal(42, answer)
	require.Equal(float32(3.14), pi)
	require.ElementsMatch([]error{err_foo, err_bar}, errs)
}
