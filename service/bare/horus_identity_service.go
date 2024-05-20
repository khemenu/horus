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
	identity "khepri.dev/horus/ent/identity"
	user "khepri.dev/horus/ent/user"
)

// IdentityService implements IdentityServiceServer
type IdentityService struct {
	client *ent.Client
	horus.UnimplementedIdentityServiceServer
}

// NewIdentityService returns a new IdentityService
func NewIdentityService(client *ent.Client) *IdentityService {
	return &IdentityService{
		client: client,
	}
}

// toProtoIdentity transforms the ent type to the pb type
func toProtoIdentity(e *ent.Identity) (*horus.Identity, error) {
	v := &horus.Identity{}
	created_date := timestamppb.New(e.CreatedDate)
	v.CreatedDate = created_date
	id := e.ID
	v.Id = id
	kind := e.Kind
	v.Kind = kind
	name := e.Name
	v.Name = name
	verifier := e.Verifier
	v.Verifier = verifier

	if edg := e.Edges.Owner; edg != nil {
		id, err := edg.ID.MarshalBinary()
		if err != nil {
			return nil, err
		}
		v.Owner = &horus.User{
			Id: id,
		}
	}
	return v, nil
}

// Create implements IdentityServiceServer.Create
func (svc *IdentityService) Create(ctx context.Context, req *horus.CreateIdentityRequest) (*horus.Identity, error) {
	identity := req.GetIdentity()
	m, err := svc.createBuilder(identity)
	if err != nil {
		return nil, err
	}
	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoIdentity(res)
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

// Get implements IdentityServiceServer.Get
func (svc *IdentityService) Get(ctx context.Context, req *horus.GetIdentityRequest) (*horus.Identity, error) {

	var (
		err error
		get *ent.Identity
	)
	id := req.GetId()
	switch req.GetView() {
	case horus.GetIdentityRequest_VIEW_UNSPECIFIED:
		fallthrough
	case horus.GetIdentityRequest_BASIC:
		get, err = svc.client.Identity.Get(ctx, id)
	case horus.GetIdentityRequest_WITH_EDGE_IDS:
		get, err = svc.client.Identity.Query().
			Where(identity.ID(id)).
			WithOwner(func(query *ent.UserQuery) {
				query.Select(user.FieldID)
			}).
			Only(ctx)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProtoIdentity(get)
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

// Update implements IdentityServiceServer.Update
func (svc *IdentityService) Update(ctx context.Context, req *horus.UpdateIdentityRequest) (*horus.Identity, error) {
	identity := req.GetIdentity()
	identityID := identity.GetId()
	m := svc.client.Identity.UpdateOneID(identityID)
	identityName := identity.GetName()
	m.SetName(identityName)
	identityVerifier := identity.GetVerifier()
	m.SetVerifier(identityVerifier)

	res, err := m.Save(ctx)
	switch {
	case err == nil:
		proto, err := toProtoIdentity(res)
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

// Delete implements IdentityServiceServer.Delete
func (svc *IdentityService) Delete(ctx context.Context, req *horus.DeleteIdentityRequest) (*emptypb.Empty, error) {
	var err error
	id := req.GetId()
	err = svc.client.Identity.DeleteOneID(id).Exec(ctx)
	switch {
	case err == nil:
		return &emptypb.Empty{}, nil
	case ent.IsNotFound(err):
		return nil, status.Errorf(codes.NotFound, "not found: %s", err)
	default:
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}

}

func (svc *IdentityService) createBuilder(identity *horus.Identity) (*ent.IdentityCreate, error) {
	m := svc.client.Identity.Create()
	identityCreatedDate := runtime.ExtractTime(identity.GetCreatedDate())
	m.SetCreatedDate(identityCreatedDate)
	identityKind := identity.GetKind()
	m.SetKind(identityKind)
	identityName := identity.GetName()
	m.SetName(identityName)
	identityVerifier := identity.GetVerifier()
	m.SetVerifier(identityVerifier)
	if identity.GetOwner() != nil {
		var identityOwner uuid.UUID
		if err := (&identityOwner).UnmarshalBinary(identity.GetOwner().GetId()); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err)
		}
		m.SetOwnerID(identityOwner)
	}
	return m, nil
}