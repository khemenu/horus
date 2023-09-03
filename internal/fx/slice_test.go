package fx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus/internal/fx"
)

func TestFilterV(t *testing.T) {
	even := fx.FilterV([]int{1, 2, 3, 4}, func(v int) bool {
		return v%2 == 0
	})
	require.Equal(t, []int{2, 4}, even)
}

func TestMapV(t *testing.T) {
	squared := fx.MapV([]int{1, 2, 3}, func(v int) int {
		return v * v
	})
	require.Equal(t, []int{1, 4, 9}, squared)
}

func TestFilterMapV(t *testing.T) {
	even := fx.FilterMapV([]int{1, 2, 3, 4}, func(v int) (int, bool) {
		return v, v%2 == 0
	})
	require.Equal(t, []int{2, 4}, even)
}

func TestAssociate(t *testing.T) {
	primes := fx.Associate([]int{2, 3, 5}, func(v int) (int, int) {
		return v, v
	})
	require.Equal(t, map[int]int{
		2: 2,
		3: 3,
		5: 5,
	}, primes)
}
