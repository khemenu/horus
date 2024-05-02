package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus/ent/proto/khepri/horus"
	"khepri.dev/horus/service/frame"
	"khepri.dev/horus/tokens"
)

const (
	TokenTypePassword = "password"
	TokenTypeAccess   = "access"
)

type TokenService struct {
	horus.UnimplementedTokenServiceServer
	*base
}

func (s *TokenService) Create(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)

	key, err := s.keyer.Key([]byte(req.Token.Id))
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	key_str := base64.RawStdEncoding.EncodeToString(key)

	// TODO: TypeBasic must be unique per user;
	// User upsert or transaction.
	// OR keep all tokens? then use only latest one?
	token, err := s.store.Token().Create(ctx, &horus.CreateTokenRequest{Token: &horus.Token{
		Id:    key_str,
		Owner: &horus.User{Id: f.Actor.ID[:]},
		Type:  tokens.TypeBasic,
	}})
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	token.Id = ""
	return token, nil
}

func (s *TokenService) Get(ctx context.Context, req *horus.GetTokenRequest) (*horus.Token, error) {
	// f := frame.Must(ctx)

	// token, err := s.client.Token.Query().
	// Where(token.And(
	// 	token.TypeEQ(tokens.TypeBasic),
	// 	token.HasOwnerWith(user.ID(f.Actor.ID)),
	// )).
	// First(ctx)
	// if err != nil {
	// 	if ent.IsNotFound(err) {
	// 		return nil, status.Error(codes.FailedPrecondition, "no password")
	// 	}

	// 	return nil, fmt.Errorf("query ")
	// }
	// s.keyer.Compare()

	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func (s *TokenService) Update(ctx context.Context, req *horus.UpdateTokenRequest) (*horus.Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (s *TokenService) Delete(ctx context.Context, req *horus.DeleteTokenRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
