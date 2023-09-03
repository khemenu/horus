package fx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus/internal/fx"
)

func TestMin(t *testing.T) {
	require := require.New(t)

	{
		v := fx.Min(3.14, 42)
		require.Equal(3.14, v)
	}

	{
		v := fx.Min(42, 3.14)
		require.Equal(3.14, v)
	}

	{
		v := fx.Min(21, 42)
		require.Equal(21, v)
	}

	{
		v := fx.Min(42, 21)
		require.Equal(21, v)
	}
}

func TestMax(t *testing.T) {
	require := require.New(t)

	{
		v := fx.Max(3.14, 42)
		require.Equal(float64(42), v)
	}

	{
		v := fx.Max(42, 3.14)
		require.Equal(float64(42), v)
	}

	{
		v := fx.Max(21, 42)
		require.Equal(42, v)
	}

	{
		v := fx.Max(42, 21)
		require.Equal(42, v)
	}
}

func TestClamp(t *testing.T) {
	require := require.New(t)

	{
		v := fx.Clamp(2, 2.718, 3.14)
		require.Equal(2.718, v)
	}

	{
		v := fx.Clamp(3, 2.718, 3.14)
		require.Equal(float64(3), v)
	}

	{
		v := fx.Clamp(4, 2.718, 3.14)
		require.Equal(3.14, v)
	}
}
