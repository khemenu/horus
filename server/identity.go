package server

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/identity"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type IdentityServiceServer struct {
	horus.UnimplementedIdentityServiceServer
	*base
}

func (s *IdentityServiceServer) Create(ctx context.Context, req *horus.CreateIdentityRequest) (*horus.Identity, error) {
	f := frame.Must(ctx)
	if req.Owner == nil {
		req.Owner = horus.UserById(f.Actor.ID)
	} else if u, err := s.isAncestorOrMeQ(ctx, f.Actor, req.Owner); err != nil {
		return nil, err
	} else if u == nil {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	return s.bare.Identity().Create(ctx, req)
}

func (s *IdentityServiceServer) Get(ctx context.Context, req *horus.GetIdentityRequest) (*horus.Identity, error) {
	f := frame.Must(ctx)

	p, err := bare.GetIdentitySpecifier(req)
	if err != nil {
		return nil, err
	}

	v, err := s.db.Identity.Query().
		Where(p).
		WithOwner(func(q *ent.UserQuery) {
			q.WithParent()
		}).
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	if v.Edges.Owner.ID == f.Actor.ID {
		// Get one of my identities.
	} else if ok, err := s.isAncestor(ctx, f.Actor, v.Edges.Owner); err != nil {
		return nil, err
	} else if !ok {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	return bare.ToProtoIdentity(v), nil
}

func (s *IdentityServiceServer) List(ctx context.Context, req *horus.ListIdentityRequest) (*horus.ListIdentityResponse, error) {
	f := frame.Must(ctx)

	u, err := s.isAncestorOrMeQ(ctx, f.Actor, req.GetOwner())
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	q := u.QueryIdentities().
		Order(identity.ByDateCreated(sql.OrderAsc()))
	if l := req.GetLimit(); l > 0 {
		q.Limit(int(l))
	}
	if t := req.GetToken(); t != nil {
		q.Where(identity.DateCreatedGT(t.AsTime()))
	}

	vs, err := q.All(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	return &horus.ListIdentityResponse{
		Items: fx.MapV(vs, func(v *ent.Identity) *horus.Identity {
			return bare.ToProtoIdentity(v)
		}),
	}, nil
}

func (s *IdentityServiceServer) Update(ctx context.Context, req *horus.UpdateIdentityRequest) (*horus.Identity, error) {
	f := frame.Must(ctx)

	p, err := bare.GetIdentitySpecifier(req.GetKey())
	if err != nil {
		return nil, err
	}

	owner, err := s.db.User.Query().
		Where(user.HasIdentitiesWith(p)).
		WithParent().
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	if ok, err := s.isAncestorOrMe(ctx, f.Actor, owner); err != nil {
		return nil, err
	} else if !ok {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	return s.bare.Identity().Update(ctx, req)
}

func (s *IdentityServiceServer) Delete(ctx context.Context, req *horus.GetIdentityRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)

	p, err := bare.GetIdentitySpecifier(req)
	if err != nil {
		return nil, err
	}

	v, err := s.db.Identity.Query().
		Where(p).
		WithOwner(func(q *ent.UserQuery) {
			q.WithParent()
		}).
		Only(ctx)
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	if ok, err := s.isAncestorOrMe(ctx, f.Actor, v.Edges.Owner); err != nil {
		return nil, err
	} else if !ok {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}

	return s.bare.Identity().Delete(ctx, req)
}
