package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/token"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/service/bare"
	"khepri.dev/horus/service/frame"
	"khepri.dev/horus/tokens"
)

type AuthService struct {
	horus.UnimplementedAuthServiceServer
	*base
}

func (s *AuthService) BasicSignIn(ctx context.Context, req *horus.BasicSignInRequest) (*horus.BasicSignInRseponse, error) {
	u, err := s.client.User.Query().
		Where(user.NameEQ(req.Username)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, fmt.Errorf("query user: %w", err)
	}

	token, err := s.client.Token.Query().
		Where(
			token.And(
				token.Type(tokens.TypeBasic),
				token.HasOwnerWith(user.ID(u.ID)),
			),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.FailedPrecondition, "user does not allow password login")
		}

		return nil, fmt.Errorf("query token: %w", err)
	}

	if key, err := base64.RawStdEncoding.DecodeString(token.Value); err != nil {
		return nil, fmt.Errorf("invalid format of basic token: %w", err)
	} else if err := tokens.Compare([]byte(req.Password), key); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "password mismatch")
	}

	ctx = frame.WithContext(ctx, &frame.Frame{Actor: u})
	access_token, err := s.bare.Token().Create(ctx, &horus.CreateTokenRequest{Token: &horus.Token{
		Type: tokens.TypeAccess,
	}})
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	return &horus.BasicSignInRseponse{Token: access_token}, nil
}

func (s *AuthService) TokenSignIn(ctx context.Context, req *horus.TokenSignInRequest) (*horus.TokenSignInResponse, error) {
	token, err := s.client.Token.Query().
		Where(
			token.And(
				token.ValueEQ(req.Token.Value),
				token.Type(tokens.TypeAccess),
			),
		).
		WithOwner().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return nil, fmt.Errorf("query token: %w", err)
	}

	return &horus.TokenSignInResponse{Token: fx.Must(bare.ToProtoToken(token))}, nil
}

func (s *AuthService) SignOut(ctx context.Context, req *horus.SingOutRequest) (*horus.SingOutResponse, error) {
	token, err := s.client.Token.Query().
		Where(
			token.And(
				token.ValueEQ(req.Token.Value),
				token.Type(tokens.TypeAccess),
			),
		).
		WithOwner().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return nil, fmt.Errorf("query token: %w", err)
	}

	ctx = frame.WithContext(ctx, &frame.Frame{
		Actor: token.Edges.Owner,
	})
	if _, err := s.service.Token().Delete(ctx, &horus.DeleteTokenRequest{
		Id: token.ID[:],
	}); err != nil {
		return nil, fmt.Errorf("delete token: %w", err)
	}

	return &horus.SingOutResponse{}, nil
}
