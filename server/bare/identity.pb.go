// Code generated by "protoc-gen-entpb". DO NOT EDIT

package bare

import (
	context "context"
	uuid "github.com/google/uuid"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	horus "khepri.dev/horus"
	ent "khepri.dev/horus/ent"
	identity "khepri.dev/horus/ent/identity"
	predicate "khepri.dev/horus/ent/predicate"
	user "khepri.dev/horus/ent/user"
)

type IdentityServiceServer struct {
	db *ent.Client
	horus.UnimplementedIdentityServiceServer
}

func NewIdentityServiceServer(db *ent.Client) *IdentityServiceServer {
	return &IdentityServiceServer{db: db}
}
func (s *IdentityServiceServer) Create(ctx context.Context, req *horus.CreateIdentityRequest) (*horus.Identity, error) {
	q := s.db.Identity.Create()
	if v := req.Name; v != nil {
		q.SetName(*v)
	}
	if v := req.Description; v != nil {
		q.SetDescription(*v)
	}
	q.SetKind(req.GetKind())
	q.SetVerifier(req.GetVerifier())
	if id, err := GetUserId(ctx, s.db, req.GetOwner()); err != nil {
		return nil, err
	} else {
		q.SetOwnerID(id)
	}

	res, err := q.Save(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return ToProtoIdentity(res), nil
}
func (s *IdentityServiceServer) Delete(ctx context.Context, req *horus.GetIdentityRequest) (*emptypb.Empty, error) {
	p, err := GetIdentitySpecifier(req)
	if err != nil {
		return nil, err
	}
	if _, err := s.db.Identity.Delete().Where(p).Exec(ctx); err != nil {
		return nil, ToStatus(err)
	}

	return &emptypb.Empty{}, nil
}
func (s *IdentityServiceServer) Get(ctx context.Context, req *horus.GetIdentityRequest) (*horus.Identity, error) {
	q := s.db.Identity.Query()
	if p, err := GetIdentitySpecifier(req); err != nil {
		return nil, err
	} else {
		q.Where(p)
	}

	res, err := QueryIdentityWithEdgeIds(q).Only(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return ToProtoIdentity(res), nil
}
func QueryIdentityWithEdgeIds(q *ent.IdentityQuery) *ent.IdentityQuery {
	q.WithOwner(func(q *ent.UserQuery) { q.Select(user.FieldID) })

	return q
}
func (s *IdentityServiceServer) Update(ctx context.Context, req *horus.UpdateIdentityRequest) (*horus.Identity, error) {
	id, err := GetIdentityId(ctx, s.db, req.GetKey())
	if err != nil {
		return nil, err
	}

	q := s.db.Identity.UpdateOneID(id)
	if v := req.Name; v != nil {
		q.SetName(*v)
	}
	if v := req.Description; v != nil {
		q.SetDescription(*v)
	}
	if v := req.Verifier; v != nil {
		q.SetVerifier(*v)
	}

	res, err := q.Save(ctx)
	if err != nil {
		return nil, ToStatus(err)
	}

	return ToProtoIdentity(res), nil
}
func ToProtoIdentity(v *ent.Identity) *horus.Identity {
	m := &horus.Identity{}
	m.Id = v.ID[:]
	m.DateCreated = timestamppb.New(v.DateCreated)
	m.Name = v.Name
	m.Description = v.Description
	m.Kind = v.Kind
	m.Verifier = v.Verifier
	if v := v.Edges.Owner; v != nil {
		m.Owner = ToProtoUser(v)
	}
	return m
}
func GetIdentityId(ctx context.Context, db *ent.Client, req *horus.GetIdentityRequest) (uuid.UUID, error) {
	var r uuid.UUID
	if v, err := uuid.FromBytes(req.GetId()); err != nil {
		return r, status.Errorf(codes.InvalidArgument, "id: %s", err)
	} else {
		return v, nil
	}
}
func GetIdentitySpecifier(req *horus.GetIdentityRequest) (predicate.Identity, error) {
	if v, err := uuid.FromBytes(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "id: %s", err)
	} else {
		return identity.IDEQ(v), nil
	}
}