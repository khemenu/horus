package horus

import (
	"context"
)

type tokenCtxKey struct{}

func Get(ctx context.Context) (*Token, bool) {
	token, ok := ctx.Value(tokenCtxKey{}).(*Token)
	return token, ok
}

func Must(ctx context.Context) *Token {
	token, ok := Get(ctx)
	if !ok || token == nil {
		panic("token not set")
	}

	return token
}

func ctxWithToken(ctx context.Context, token *Token) context.Context {
	return context.WithValue(ctx, tokenCtxKey{}, token)
}
