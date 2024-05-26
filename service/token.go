package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"khepri.dev/horus"
	"khepri.dev/horus/service/frame"
)

type TokenService struct {
	horus.UnimplementedTokenServiceServer
	*base
}

func (s *TokenService) Create(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	switch req.Token.Type {
	case horus.TokenTypeBasic:
		return s.createBasic(ctx, req)
	case horus.TokenTypeRefresh:
		return s.createAccess(ctx, req)
	case horus.TokenTypeAccess:
		return s.createAccess(ctx, req)
	default:
		break
	}

	f := frame.Must(ctx)
	req.Token.Owner = &horus.User{Id: f.Actor.ID[:]}
	return s.bare.Token().Create(ctx, req)
}

func (s *TokenService) createBasic(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)

	key, err := s.keyer.Key([]byte(req.Token.Value))
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	key_str := ""
	if b, err := proto.Marshal(key); err != nil {
		return nil, fmt.Errorf("marshal key: %w", err)
	} else {
		key_str = base64.RawStdEncoding.EncodeToString(b)
	}

	// TODO: TypeBasic must be unique per user;
	// Use upsert or transaction.
	// OR keep all tokens? then use only latest one?
	token, err := s.bare.Token().Create(ctx, &horus.CreateTokenRequest{Token: &horus.Token{
		Value: key_str,
		Owner: &horus.User{Id: f.Actor.ID[:]},
		Type:  horus.TokenTypeBasic,
	}})
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	token.Value = ""
	return token, nil
}

func (s *TokenService) createRefresh(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)
	v, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return s.bare.Token().Create(ctx, &horus.CreateTokenRequest{
		Token: &horus.Token{
			Value:       v,
			Type:        horus.TokenTypeRefresh,
			DateExpired: timestamppb.New(time.Now().Add(24 * time.Hour)),
			Owner:       &horus.User{Id: f.Actor.ID[:]},
		},
	})
}

func (s *TokenService) createAccess(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)
	v, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return s.bare.Token().Create(ctx, &horus.CreateTokenRequest{
		Token: &horus.Token{
			Value:       v,
			Type:        horus.TokenTypeAccess,
			DateExpired: timestamppb.New(time.Now().Add(24 * time.Hour)),
			Owner:       &horus.User{Id: f.Actor.ID[:]},
			Parent:      req.Token.GetParent(),
		},
	})
}

func (*TokenService) generateToken() (string, error) {
	const Size = 128
	charSet := []rune("-.ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz0123456789~")

	rst := make([]rune, Size)
	buff := make([]byte, 8)

	for i := range rst {
		if _, err := rand.Read(buff); err != nil {
			return "", fmt.Errorf("crypto rand: %w", err)
		}

		idx := binary.LittleEndian.Uint64(buff) % uint64(len(charSet))
		rst[i] = charSet[idx]
	}

	return string(rst), nil
}

func (s *TokenService) Get(ctx context.Context, req *horus.GetTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)

	req.View = horus.GetTokenRequest_WITH_EDGE_IDS
	token, err := s.bare.Token().Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(token.Id, f.Actor.ID[:]) {
		return nil, status.Error(codes.NotFound, "not found")
	}

	switch token.Type {
	default:
		token.Value = ""
	}

	return token, nil
}

func (s *TokenService) Update(ctx context.Context, req *horus.UpdateTokenRequest) (*horus.Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *TokenService) Delete(ctx context.Context, req *horus.DeleteTokenRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)
	token, err := s.bare.Token().Get(ctx, &horus.GetTokenRequest{Id: req.Id})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "not found")
		}

		return nil, fmt.Errorf("get token: %w", err)
	}
	if !bytes.Equal(token.GetOwner().Id, f.Actor.ID[:]) {
		return nil, status.Error(codes.NotFound, "not found")
	}

	token.DateExpired = timestamppb.Now()
	if _, err := s.bare.Token().Update(ctx, &horus.UpdateTokenRequest{Token: token}); err != nil {
		return nil, fmt.Errorf("update token expired date: %w", err)
	}

	return &emptypb.Empty{}, nil
}
