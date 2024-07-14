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

func TestCond(t *testing.T) {
	{
		v := fx.Cond(true, 36, 42)
		require.Equal(t, v, 36)
	}

	{
		v := fx.Cond(false, 36, 42)
		require.Equal(t, v, 42)
	}
}

func TestAnd(t *testing.T) {
	o := true
	x := false

	for _, vs := range [][]bool{
		{o},

		{x, x},
		{o, o},

		{x, x, x},
		{x, o, x},
		{x, x, o},
		{o, o, o},

		{x, x, x, x},
		{x, o, x, x},
		{x, x, o, x},
		{x, o, o, x},
		{x, x, x, o},
		{x, o, x, o},
		{x, x, o, o},
		{o, o, o, o},
	} {
		v := fx.And(vs[1:]...)
		require.Equal(t, vs[0], v)
	}
}

func TestOr(t *testing.T) {
	o := true
	x := false

	for _, vs := range [][]bool{
		{x},

		{x, x},
		{o, o},

		{x, x, x},
		{o, o, x},
		{o, x, o},
		{o, o, o},

		{x, x, x, x},
		{o, o, x, x},
		{o, x, o, x},
		{o, o, o, x},
		{o, x, x, o},
		{o, o, x, o},
		{o, x, o, o},
		{o, o, o, o},
	} {
		v := fx.Or(vs[1:]...)
		require.Equal(t, vs[0], v)
	}
}
