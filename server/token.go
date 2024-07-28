package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
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
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type TokenServiceServer struct {
	horus.UnimplementedTokenServiceServer
	*base
}

func (s *TokenServiceServer) hasPermission(ctx context.Context, v *ent.Token) error {
	if v == nil {
		return nil
	}
	if v.Type != horus.TokenTypeAccess {
		return nil
	}

	n, err := s.db.Token.Query().
		Where(
			token.IDEQ(v.ID),
			token.HasParentWith(token.TypeNEQ(horus.TokenTypePassword)),
		).
		Count(ctx)
	if err != nil {
		return bare.ToStatus(err)
	}
	if n > 0 {
		return status.Error(codes.PermissionDenied, "it is not allowed to access a token service using a token created by another access token or a refresh token")
	}

	return nil
}

func (s *TokenServiceServer) Create(ctx context.Context, req *horus.CreateTokenRequest) (v *horus.Token, err error) {
	f := frame.Must(ctx)
	if err := s.hasPermission(ctx, f.Token); err != nil {
		return nil, err
	}
	if req == nil {
		req = &horus.CreateTokenRequest{}
	}

	var owner_id uuid.UUID
	if owner := req.Owner; owner == nil {
		owner_id = f.Actor.ID
	} else {
		p, err := bare.GetUserSpecifier(owner)
		if err != nil {
			return nil, err
		}

		id, err := s.db.User.Query().Where(p, user.HasParentWith(user.IDEQ(f.Actor.ID))).OnlyID(ctx)
		if err != nil {
			return nil, bare.ToStatus(err)
		}

		owner_id = id
	}

	req.Owner = horus.UserById(owner_id)
	if f.Token != nil {
		req.Parent = horus.TokenById(f.Token.ID)
	}

	switch req.GetType() {
	case horus.TokenTypePassword:
		v, err = s.createBasic(ctx, req)
	case horus.TokenTypeRefresh:
		v, err = s.createBearer(ctx, req, horus.TokenTypeRefresh)
	case horus.TokenTypeAccess:
		v, err = s.createBearer(ctx, req, horus.TokenTypeAccess)
	default:
		return nil, status.Error(codes.Unimplemented, "unimplemented")
	}
	if err != nil {
		return
	}

	v.Owner = &horus.User{Id: owner_id[:]}
	return
}

func (s *TokenServiceServer) createBasic(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)

	pw := req.GetValue()
	if pw == "" {
		return nil, status.Error(codes.InvalidArgument, `"Token.value" must be provided`)
	}

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

	p, err := bare.GetUserSpecifier(req.GetOwner())
	if err != nil {
		return nil, err
	}

	owner, err := s.db.User.Query().Where(p).Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	v, err := entutils.WithTxV(ctx, s.db, func(tx *ent.Tx) (*ent.Token, error) {
		_, err := tx.Token.Delete().
			Where(
				token.TypeEQ(horus.TokenTypePassword),
				token.HasOwnerWith(user.IDEQ(owner.ID)),
			).
			Exec(ctx)
		if err != nil {
			return nil, fmt.Errorf("delete existing password: %w", err)
		}

		// TODO: use bare server to create?
		q := tx.Token.Create().
			SetValue(key_str).
			SetType(horus.TokenTypePassword).
			SetUseCountLimit(req.GetUseCountLimit()).
			SetDateExpired(time.Now().Add(10 * 365 * 24 * time.Hour)).
			SetOwnerID(owner.ID)
		if f.Token != nil {
			q.SetParentID(f.Token.ID)
		}

		now := time.Now()
		if owner.DateUnlocked == nil || owner.DateUnlocked.After(now) {
			err := tx.User.UpdateOneID(owner.ID).
				SetDateUnlocked(now).
				Exec(ctx)
			if err != nil {
				return nil, err
			}
		}

		v, err := q.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("create token: %w", err)
		}

		return v, nil
	})
	if err != nil {
		return nil, err
	}

	v.Value = ""
	return bare.ToProtoToken(v), nil
}

func (s *TokenServiceServer) createBearer(ctx context.Context, req *horus.CreateTokenRequest, t string) (*horus.Token, error) {
	const TokenLength = 128

	if req.GetValue() != "" {
		return nil, status.Error(codes.InvalidArgument, "value ot bearer cannot be set manually")
	}

	var token [TokenLength]byte
	_, err := io.ReadFull(rand.Reader, token[:])
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	key, err := s.keyer.Key(token[:])
	if err != nil {
		return nil, fmt.Errorf("key: %w", err)
	}

	key_str := ""
	if b, err := proto.Marshal(key); err != nil {
		return nil, fmt.Errorf("marshal key: %w", err)
	} else {
		key_str = base64.RawStdEncoding.EncodeToString(b)
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

	owner_id, err := bare.GetUserId(ctx, s.db, req.GetOwner())
	if err != nil {
		return nil, err
	}

	v, err := s.bare.Token().Create(ctx, &horus.CreateTokenRequest{
		Value:  key_str,
		Type:   t,
		Owner:  horus.UserById(owner_id),
		Parent: req.GetParent(),

		DateExpired: ts_expired,
	})
	if err != nil {
		return nil, err
	}

	v.Value = base64.RawStdEncoding.EncodeToString(append(v.Id, token[:]...))
	return v, nil
}

func (s *TokenServiceServer) Get(ctx context.Context, req *horus.GetTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)

	p, err := bare.GetTokenSpecifier(req)
	if err != nil {
		return nil, err
	}

	v, err := bare.QueryTokenWithEdgeIds(s.db.Token.Query()).
		Where(p,
			token.HasOwnerWith(user.IDEQ(f.Actor.ID)),
			token.DateExpiredGT(time.Now()),
		).
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	} else {
		v.Value = ""
	}

	return bare.ToProtoToken(v), nil
}

func (s *TokenServiceServer) List(ctx context.Context, req *horus.ListTokenRequest) (*horus.ListTokenResponse, error) {
	f := frame.Must(ctx)

	q := s.db.Token.Query().
		Order(token.ByDateCreated(sql.OrderDesc())).
		Where(
			token.HasOwnerWith(user.IDEQ(f.Actor.ID)),
			token.DateExpiredGT(time.Now()),
		)
	if l := req.GetLimit(); l > 0 {
		q.Limit(int(l))
	}
	if t := req.GetToken(); t != nil {
		q.Where(token.DateCreatedLT(t.AsTime()))
	}

	var (
		vs  []*ent.Token
		err error
	)
	switch k := req.GetKey().(type) {
	case *horus.ListTokenRequest_Type:
		vs, err = q.Where(token.TypeEQ(k.Type)).All(ctx)
	}
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	return &horus.ListTokenResponse{
		Items: fx.MapV(vs, func(v *ent.Token) *horus.Token {
			v.Value = ""
			return bare.ToProtoToken(v)
		}),
	}, nil
}

func (s *TokenServiceServer) Update(ctx context.Context, req *horus.UpdateTokenRequest) (*horus.Token, error) {
	f := frame.Must(ctx)
	if err := s.hasPermission(ctx, f.Token); err != nil {
		return nil, err
	}

	v, err := s.Get(ctx, req.GetKey())
	if err != nil {
		return nil, err
	}

	req.Key = horus.TokenByIdV(v.Id)
	v, err = s.bare.Token().Update(ctx, req)
	if err != nil {
		return nil, err
	}

	v.Value = ""
	return v, nil
}

func (s *TokenServiceServer) Delete(ctx context.Context, req *horus.GetTokenRequest) (*emptypb.Empty, error) {
	_, err := s.Update(ctx, &horus.UpdateTokenRequest{
		Key:         req,
		DateExpired: timestamppb.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
