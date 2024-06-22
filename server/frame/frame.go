package frame

import (
	"context"

	"khepri.dev/horus/ent"
)

type Frame struct {
	Actor         *ent.User
	Token         *ent.Token
	ActingAccount *ent.Account
}

func (f *Frame) MustGetActingAccount() *ent.Account {
	if f.ActingAccount == nil {
		panic("after Account.Get, if the target wasn't myself, acting account must be set")
	}

	return f.ActingAccount
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
