package fx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus/internal/fx"
)

func TestDefault(t *testing.T) {
	v := 0
	fx.Default(&v, 42)
	require.Equal(t, 42, v)
}

func TestFallback(t *testing.T) {
	a := 0
	b := 42
	{
		rst := fx.Fallback(a, b)
		require.Equal(t, b, rst)
	}

	a = 36
	{
		rst := fx.Fallback(a, b)
		require.Equal(t, a, rst)
	}
}
