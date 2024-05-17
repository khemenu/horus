// Code generated by protoc-gen-entgrpc. DO NOT EDIT.
package bare

import (
	context "context"
	runtime "entgo.io/contrib/entproto/runtime"
	sqlgraph "entgo.io/ent/dialect/sql/sqlgraph"
	uuid "github.com/google/uuid"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	horus "khepri.dev/horus"
	ent "khepri.dev/horus/ent"
	token "khepri.dev/horus/ent/token"
	user "khepri.dev/horus/ent/user"
)

// TokenService implements TokenServiceServer
type TokenService struct {
	client *ent.Client
	horus.UnimplementedTokenServiceServer
}

// NewTokenService returns a new TokenService
func NewTokenService(client *ent.Client) *TokenService {
	return &TokenService{
		client: client,
	}
}

// toProtoToken transforms the ent type to the pb type
func toProtoToken(e *ent.Token) (*horus.Token, error) {
	v := &horus.Token{}
	date_created := timestamppb.New(e.DateCreated)
	v.DateCreated = date_created
	date_expired := timestamppb.New(e.DateExpired)
	v.DateExpired = date_expired
	id, err := e.ID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	v.Id = id
	name := e.Name
	v.Name = name
	_type := e.Type
	v.Type = _type
	value := e.Value
	v.Value = value

	if edg := e.Edges.Owner; edg != nil {
		id, err := edg.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		v.Owner = &horus.User{
			Id: id,
		}
	}

	if edg := e.Edges.Parent; edg != nil {
		id, err := edg.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		v.Parent = &horus.Token{
			Id: id,
		}
	}
	return v, nil
}

// Create implements TokenServiceServer.Create
func (svc *TokenService) Create(ctx context.Context, req *horus.CreateTokenRequest) (*horus.Token, error) {
	token := req.GetToken()
	m, err := svc.createBuilder(token)
	if err != nil {
		return nil, err
	}
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoToken(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Get implements TokenServiceServer.Get
func (svc *TokenService) Get(ctx context.Context, req *horus.GetTokenRequest) (*horus.Token, error) {

	var (
		err error
		get *ent.Token
	)
	var id uuid.UUID
	if err := (&id).UnmarshalBinary(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}
	switch req.GetView() {
	case horus.GetTokenRequest_VIEW_UNSPECIFIED:
		fallthrough
	case horus.GetTokenRequest_BASIC:
		get, err = svc.client.Token.Get(ctx, id)
	case horus.GetTokenRequest_WITH_EDGE_IDS:
		get, err = svc.client.Token.Query().
			Where(token.ID(id)).
			WithOwner(func(query *ent.UserQuery) {
				query.Select(user.FieldID)
			}).
			WithParent(func(query *ent.TokenQuery) {
				query.Select(token.FieldID)
			}).
			Only(ctx)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProtoToken(get)
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Update implements TokenServiceServer.Update
func (svc *TokenService) Update(ctx context.Context, req *horus.UpdateTokenRequest) (*horus.Token, error) {
	token := req.GetToken()
	var tokenID uuid.UUID
	if err := (&tokenID).UnmarshalBinary(token.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}
	m := svc.client.Token.UpdateOneID(tokenID)
	tokenDateExpired := runtime.ExtractTime(token.GetDateExpired())
	m.SetDateExpired(tokenDateExpired)
	tokenName := token.GetName()
	m.SetName(tokenName)

	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoToken(res)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal error: %s", err)
		}
		return proto, nil
	case sqlgraph.IsUniqueConstraintError(err):
		return nil, status.Errorf(codes.AlreadyExists, "already exists: %s", err)
	case ent.IsConstraintError(err):
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Delete implements TokenServiceServer.Delete
func (svc *TokenService) Delete(ctx context.Context, req *horus.DeleteTokenRequest) (*emptypb.Empty, error) {
	var err error
	var id uuid.UUID
	if err := (&id).UnmarshalBinary(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
	}
	err = svc.client.Token.DeleteOneID(id).Exec(ctx)
	switch {
	case err == nil:
		return &emptypb.Empty{}, nil
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

func (svc *TokenService) createBuilder(token *horus.Token) (*ent.TokenCreate, error) {
	m := svc.client.Token.Create()
	tokenDateCreated := runtime.ExtractTime(token.GetDateCreated())
	m.SetDateCreated(tokenDateCreated)
	tokenDateExpired := runtime.ExtractTime(token.GetDateExpired())
	m.SetDateExpired(tokenDateExpired)
	tokenName := token.GetName()
	m.SetName(tokenName)
	tokenType := token.GetType()
	m.SetType(tokenType)
	tokenValue := token.GetValue()
	m.SetValue(tokenValue)
	if token.GetOwner() != nil {
		var tokenOwner uuid.UUID
		if err := (&tokenOwner).UnmarshalBinary(token.GetOwner().GetId()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.SetOwnerID(tokenOwner)
	}
	if token.GetParent() != nil {
		var tokenParent uuid.UUID
		if err := (&tokenParent).UnmarshalBinary(token.GetParent().GetId()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.SetParentID(tokenParent)
	}
	return m, nil
}
