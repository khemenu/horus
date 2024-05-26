package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/token"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/service/bare"
	"khepri.dev/horus/service/frame"
	"khepri.dev/horus/tokens"
)

type AuthService struct {
	horus.UnimplementedAuthServiceServer
	*base
}

func (s *AuthService) BasicSignIn(ctx context.Context, req *horus.BasicSignInRequest) (*horus.BasicSignInResponse, error) {
	u, err := s.client.User.Query().
		Where(user.AliasEQ(req.Username)).
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
				token.Type(horus.TokenTypeBasic),
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
	access_token, err := s.service.Token().Create(ctx, &horus.CreateTokenRequest{Token: &horus.Token{
		Type: horus.TokenTypeAccess,
	}})
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	return &horus.BasicSignInResponse{
		Token: access_token,
	}, nil
}

func (s *AuthService) TokenSignIn(ctx context.Context, req *horus.TokenSignInRequest) (*horus.TokenSignInResponse, error) {
	token, err := s.client.Token.Query().
		Where(
			token.ValueEQ(req.Token),
			token.Type(horus.TokenTypeAccess),
		).
		WithOwner().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return nil, fmt.Errorf("query token: %w", err)
	}

	return &horus.TokenSignInResponse{
		Token: fx.Must(bare.ToProtoToken(token)),
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, req *horus.RefreshRequest) (*horus.RefreshResponse, error) {
	return entutils.WithTxV(ctx, s.client, func(tx *ent.Tx) (*horus.RefreshResponse, error) {
		refresh_token, err := tx.Token.Query().
			Where(token.And(
				token.ValueEQ(req.Token),
				token.Type(horus.TokenTypeRefresh),
				token.DateExpiredGT(time.Now()),
			)).
			Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, status.Error(codes.Unauthenticated, "invalid token")
			}

			return nil, fmt.Errorf("query token: %w", err)
		}

		now := time.Now()
		if _, err := tx.Token.Update().
			Where(token.And(
				token.Type(horus.TokenTypeAccess),
				token.HasParentWith(token.ID(refresh_token.ID)),
				token.DateExpiredGT(now),
			)).
			SetDateExpired(now).
			Save(ctx); err != nil {
			if ent.IsNotFound(err) {
				return nil, status.Error(codes.Unauthenticated, "expires previous tokens")
			}

			return nil, fmt.Errorf("update token: %w", err)
		}

		ctx = frame.WithContext(ctx, &frame.Frame{Actor: refresh_token.Edges.Owner})
		access_token, err := (&TokenService{base: s.withClient(tx.Client())}).Create(ctx, &horus.CreateTokenRequest{Token: &horus.Token{
			Type:   horus.TokenTypeAccess,
			Parent: &horus.Token{Id: refresh_token.ID[:]},
		}})
		if err != nil {
			return nil, fmt.Errorf("create access token: %w", err)
		}

		return &horus.RefreshResponse{
			Token: access_token,
		}, nil
	})
}

func (s *AuthService) VerifyOtp(ctx context.Context, req *horus.VerifyOtpRequest) (*horus.VerifyOtpResponse, error) {
	now := time.Now()
	n, err := s.client.Token.Update().
		Where(token.And(
			token.ValueEQ(req.Value),
			token.Type(horus.TokenTypeOtp),
			token.DateCreatedGT(now),
		)).
		SetDateExpired(now).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("update token: %w", err)
	}
	if n != 0 {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return &horus.VerifyOtpResponse{}, nil
}

func (s *AuthService) SignOut(ctx context.Context, req *horus.SingOutRequest) (*horus.SingOutResponse, error) {
	token, err := s.client.Token.Query().
		Where(
			token.And(
				token.ValueEQ(req.Token),
				token.Type(horus.TokenTypeAccess),
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
