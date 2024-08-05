package horus

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthUnaryInterceptor(signIn func(ctx context.Context, in *TokenSignInRequest) (*TokenSignInResponse, error)) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		token := ""
		if md, ok := metadata.FromIncomingContext(ctx); !ok {
			return nil, status.Error(codes.InvalidArgument, "missing metadata")
		} else if entries := md.Get("authorization"); len(entries) > 0 && strings.HasPrefix(entries[0], "Bearer ") {
			token, _ = strings.CutPrefix(entries[0], "Bearer ")
		} else if entries := md.Get(TokenKeyName); len(entries) > 0 {
			token = entries[0]
		} else if entries := md.Get("cookie"); len(entries) > 0 {
			prefix := fmt.Sprintf("%s=", TokenKeyName)
			for _, cookie := range strings.Split(entries[0], "; ") {
				if !strings.HasPrefix(cookie, prefix) {
					continue
				}

				kv := strings.SplitN(cookie, "=", 2)
				if len(kv) != 2 {
					continue
				}

				token = kv[1]
				break
			}
		}
		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "no access token")
		}

		res, err := signIn(metadata.NewOutgoingContext(ctx, metadata.MD{}), &TokenSignInRequest{
			Token: token,
		})
		if err != nil {
			return nil, err
		}

		ctx = ctxWithToken(ctx, res.Token)
		ctx = metadata.AppendToOutgoingContext(ctx, TokenKeyName, res.Token.Value)
		return handler(ctx, req)
	}
}
