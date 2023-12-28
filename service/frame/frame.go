package frame

import (
	"context"

	"khepri.dev/horus/ent"
)

type Frame struct {
	Actor         *ent.User
	ActingAccount *ent.Account
}

var key = struct{}{}

func New() *Frame {
	return &Frame{}
}

func Get(ctx context.Context) (*Frame, bool) {
	frame, ok := ctx.Value(key).(*Frame)
	return frame, ok
}

func Must(ctx context.Context) *Frame {
	frame, ok := Get(ctx)
	if !ok || frame == nil {
		panic("frame must be set")
	}

	return frame
}

func WithContext(ctx context.Context, value *Frame) context.Context {
	return context.WithValue(ctx, key, value)
}
