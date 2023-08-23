package frame

import "context"

type frameCtxKey struct{}

func FromCtx(ctx context.Context) (Frame, bool) {
	f, ok := ctx.Value(frameCtxKey{}).(Frame)
	if !ok {
		return nil, false
	}

	return f, ok
}

func MustFromCtx(ctx context.Context) Frame {
	f, ok := FromCtx(ctx)
	if !ok {
		panic("no frame in the context")
	}

	return f
}

func WithCtx(ctx context.Context, frame Frame) context.Context {
	return context.WithValue(ctx, frameCtxKey{}, frame)
}
