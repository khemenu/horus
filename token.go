package horus

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type tokenCtxKey struct{}

func From(ctx context.Context) (*Token, bool) {
	token, ok := ctx.Value(tokenCtxKey{}).(*Token)
	return token, ok
}

func Must(ctx context.Context) *Token {
	token, ok := From(ctx)
	if !ok || token == nil {
		panic("token not set")
	}

	return token
}

func ctxWithToken(ctx context.Context, token *Token) context.Context {
	return context.WithValue(ctx, tokenCtxKey{}, token)
}

func WithToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, TokenKeyName, token)
}
