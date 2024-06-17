package server

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/token"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type TokenServiceServer struct {
	horus.UnimplementedTokenServiceServer
	*base
}

func (s *TokenServiceServer) Create(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)
	if req != nil {
		req.Parent = horus.TokenById(f.Token.ID)
	}

	switch req.GetType() {
	case horus.TokenTypeBasic:
		return s.createBasic(ctx, req)
	case horus.TokenTypeRefresh:
		return s.createBearer(ctx, req, horus.TokenTypeRefresh)
	case horus.TokenTypeAccess:
		return s.createBearer(ctx, req, horus.TokenTypeAccess)
	default:
		return nil, status.Error(codes.Unimplemented, "unimplemented")
	}
}

func (s *TokenServiceServer) createBasic(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)

	pw := req.GetValue()
	if pw == "" {
		return nil, status.Error(codes.InvalidArgument, `"Token.value" must be provided`)
	}
	pw = strings.TrimSpace(pw)

	key, err := s.keyer.Key([]byte(pw))
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	key_str := ""
	if b, err := proto.Marshal(key); err != nil {
		return nil, fmt.Errorf("marshal key: %w", err)
	} else {
		key_str = base64.RawStdEncoding.EncodeToString(b)
	}

	v, err := entutils.WithTxV(ctx, s.db, func(tx *ent.Tx) (*ent.Token, error) {
		_, err := tx.Token.Delete().
			Where(
				token.TypeEQ(horus.TokenTypeBasic),
				token.HasOwnerWith(user.IDEQ(f.Actor.ID)),
			).
			Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("delete existing basic tokens: %w", err)
		}

		v, err := tx.Token.Create().
			SetValue(key_str).
			SetType(horus.TokenTypeBasic).
			SetDateExpired(time.Now().Add(10 * 365 * 24 * time.Hour)).
			SetOwnerID(f.Actor.ID).
			SetParentID(f.Token.ID).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("create token: %w", err)
		}

		return v, nil
	})
	if err != nil {
		return nil, err
	}

	return &horus.Token{
		Id:          v.ID[:],
		Type:        horus.TokenTypeBasic,
		DateExpired: timestamppb.New(v.DateExpired),
		DateCreated: timestamppb.New(v.DateCreated),
	}, nil
}

func (s *TokenServiceServer) createBearer(ctx context.Context, req *horus.CreateTokenRequest, t string) (*horus.Token, error) {
	f := frame.Must(ctx)
	v, err := s.generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	var date_expired time.Time
	switch t {
	case horus.TokenTypeRefresh:
		date_expired = time.Now().Add(10 * 365 * 24 * time.Hour)
	case horus.TokenTypeAccess:
		date_expired = time.Now().Add(24 * time.Hour)

	default:
		panic("invalid bearer token type")
	}

	var ts_expired *timestamppb.Timestamp
	if d := req.GetDateExpired(); d != nil {
		ts_expired = d
	} else {
		ts_expired = timestamppb.New(date_expired)
	}

	owner_id := f.Actor.ID[:]
	if get_child := req.GetOwner(); get_child != nil {
		if bytes.Equal(get_child.GetId(), owner_id) || get_child.GetAlias() == f.Actor.Alias {
			// get_child was not a child but myself.
		} else {
			p, err := bare.GetUserSpecifier(get_child)
			if err != nil {
				return nil, err
			}

			child, err := s.db.User.Query().
				Where(p, user.HasParentWith(user.IDEQ(f.Actor.ID))).
				OnlyID(ctx)
			if err != nil {
				return nil, bare.ToStatus(err)
			}

			owner_id = child[:]
		}
	}

	return s.bare.Token().Create(ctx, &horus.CreateTokenRequest{
		Value:  v,
		Type:   t,
		Owner:  horus.UserByIdV(owner_id),
		Parent: req.GetParent(),

		DateExpired: ts_expired,
	})
}

func (*TokenServiceServer) generateToken() (string, error) {
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

func (s *TokenServiceServer) Get(ctx context.Context, req *horus.GetTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)

	token, err := s.bare.Token().Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(token.Owner.Id, f.Actor.ID[:]) {
		return nil, status.Error(codes.NotFound, "not found")
	}

	token.Value = ""
	return token, nil
}

func (s *TokenServiceServer) Update(ctx context.Context, req *horus.UpdateTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)

	p, err := bare.GetTokenSpecifier(req.GetKey())
	if err != nil {
		return nil, err
	}

	id, err := s.db.Token.Query().
		Where(p, token.HasOwnerWith(user.IDEQ(f.Actor.ID))).
		OnlyID(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	req.Key = horus.TokenById(id)
	return s.bare.Token().Update(ctx, req)
}

func (s *TokenServiceServer) Delete(ctx context.Context, req *horus.GetTokenRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)

	token, err := s.bare.Token().Get(ctx, &horus.GetTokenRequest{Key: &horus.GetTokenRequest_Id{
		Id: req.GetId(),
	}})
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	if !bytes.Equal(token.GetOwner().Id, f.Actor.ID[:]) {
		return nil, status.Error(codes.NotFound, "not found")
	}

	if _, err := s.bare.Token().Update(ctx, &horus.UpdateTokenRequest{
		Key:         horus.TokenByIdV(token.Id),
		DateExpired: timestamppb.Now(),
	}); err != nil {
		return nil, fmt.Errorf("update token expired date: %w", err)
	}

	return &emptypb.Empty{}, nil
}
