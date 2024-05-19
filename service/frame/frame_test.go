package frame_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"khepri.dev/horus/service/frame"
)

func TestFrame(t *testing.T) {
	require := require.New(t)

	ctx := context.Background()
	_, ok := frame.Get(ctx)
	require.False(ok)

	ctx = frame.WithContext(ctx, nil)
	require.Panics(func() { frame.Must(ctx) })

	frame_origin := frame.New()
	ctx = frame.WithContext(ctx, frame_origin)
	frame_retrieved, ok := frame.Get(ctx)
	require.True(ok)
	require.Same(frame_origin, frame_retrieved)
	require.NotPanics(func() { frame.Must(ctx) })
}
