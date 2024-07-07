package server

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/token"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
	"khepri.dev/horus/tokens"
)

type AuthService struct {
	horus.UnimplementedAuthServiceServer
	*base
}

func (s *AuthService) BasicSignIn(ctx context.Context, req *horus.BasicSignInRequest) (*horus.BasicSignInResponse, error) {
	token, err := s.db.Token.Query().
		Where(
			token.Type(horus.TokenTypePassword),
			token.HasOwnerWith(user.AliasEQ(req.GetUsername())),
		).
		WithOwner().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.Unauthenticated, "")
		}

		return nil, fmt.Errorf("query token: %w", err)
	}

	if key, err := base64.RawStdEncoding.DecodeString(token.Value); err != nil {
		return nil, fmt.Errorf("invalid format of basic token: %w", err)
	} else if err := tokens.Compare([]byte(req.Password), key); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "")
	}

	ctx = frame.WithContext(ctx, &frame.Frame{
		Actor: token.Edges.Owner,
		Token: token,
	})
	v, err := s.covered.Token().Create(ctx, &horus.CreateTokenRequest{
		Type: horus.TokenTypeAccess,
	})
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	return &horus.BasicSignInResponse{
		Token: v,
	}, nil
}

func (s *AuthService) verifyToken(ctx context.Context, token_str string, token_type string) (*ent.Token, error) {
	t, err := base64.RawStdEncoding.DecodeString(token_str)
	if err == nil && len(t) < 16 {
		err = fmt.Errorf("too short")
	}
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid token format: %s", err.Error())
	}

	id, err := uuid.FromBytes(t[:16])
	if err != nil {
		panic(fmt.Errorf("too short: %w", err))
	}

	v, err := s.db.Token.Query().
		Where(
			token.IDEQ(id),
			token.TypeEQ(token_type),
			token.DateExpiredGT(time.Now()),
		).
		WithOwner().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.Unauthenticated, "")
		}

		return nil, fmt.Errorf("query token: %w", err)
	}

	if key, err := base64.RawStdEncoding.DecodeString(v.Value); err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	} else if err := tokens.Compare(t[16:], key); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "")
	}

	return v, nil
}

func (s *AuthService) TokenSignIn(ctx context.Context, req *horus.TokenSignInRequest) (*horus.TokenSignInResponse, error) {
	v, err := s.verifyToken(ctx, req.GetToken(), horus.TokenTypeAccess)
	if err != nil {
		return nil, err
	}

	if frame, ok := frame.Get(ctx); ok {
		frame.Actor = v.Edges.Owner
		frame.Token = v
	}

	return &horus.TokenSignInResponse{
		Token: bare.ToProtoToken(v),
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, req *horus.RefreshRequest) (*horus.RefreshResponse, error) {
	refresh_token, err := s.verifyToken(ctx, req.GetToken(), horus.TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	return entutils.WithTxV(ctx, s.db, func(tx *ent.Tx) (*horus.RefreshResponse, error) {
		now := time.Now()
		if _, err := tx.Token.Update().
			Where(
				token.Type(horus.TokenTypeAccess),
				token.HasParentWith(token.ID(refresh_token.ID)),
				token.DateExpiredGT(now),
			).
			SetDateExpired(now).
			Save(ctx); err != nil {
			if !ent.IsNotFound(err) {
				return nil, fmt.Errorf("update token: %w", err)
			}
		}

		ctx = frame.WithContext(ctx, &frame.Frame{
			Actor: refresh_token.Edges.Owner,
			Token: refresh_token,
		})
		access_token, err := (&TokenServiceServer{base: s.withClient(tx.Client())}).Create(ctx, &horus.CreateTokenRequest{
			Type: horus.TokenTypeAccess,
		})
		if err != nil {
			return nil, err
		}

		return &horus.RefreshResponse{
			Token: access_token,
		}, nil
	})
}

func (s *AuthService) VerifyOtp(ctx context.Context, req *horus.VerifyOtpRequest) (*horus.VerifyOtpResponse, error) {
	now := time.Now()
	n, err := s.db.Token.Update().
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
	token, err := s.db.Token.Query().
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
	if _, err := s.covered.Token().Delete(ctx, horus.TokenById(token.ID)); err != nil {
		return nil, fmt.Errorf("delete token: %w", err)
	}

	return &horus.SingOutResponse{}, nil
}
