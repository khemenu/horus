package fx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus/internal/fx"
)

func TestMapV(t *testing.T) {
	squared := fx.MapV([]int{1, 2, 3}, func(v int) int {
		return v * 2
	})
	require.Equal(t, []int{2, 4, 6}, squared)
}
