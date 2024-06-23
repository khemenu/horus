package server

import (
	"bytes"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/server/frame"
)

type UserServiceServer struct {
	horus.UnimplementedUserServiceServer
	*base
}

func (s *UserServiceServer) Create(ctx context.Context, req *horus.CreateUserRequest) (*horus.User, error) {
	f := frame.Must(ctx)
	if req == nil {
		req = &horus.CreateUserRequest{}
	}
	if req.Parent != nil {
		return nil, status.Error(codes.InvalidArgument, "parent cannot be set manually")
	}
	req.Parent = horus.UserById(f.Actor.ID)
	return s.bare.User().Create(ctx, req)
}

func (s *UserServiceServer) Get(ctx context.Context, req *horus.GetUserRequest) (*horus.User, error) {
	f := frame.Must(ctx)
	if req.GetKey() == nil || req.GetAlias() == "_me" {
		req = horus.UserById(f.Actor.ID)
	}

	v, err := s.bare.User().Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(f.Actor.ID[:], v.Id) {
		return v, nil
	}
	if bytes.Equal(f.Actor.ID[:], v.GetParent().GetId()) {
		return v, nil
	}

	return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
}

func (s *UserServiceServer) Update(ctx context.Context, req *horus.UpdateUserRequest) (*horus.User, error) {
	u, err := s.Get(ctx, req.GetKey())
	if err != nil {
		return nil, err
	}
	if req != nil {
		req.Parent = nil
	}

	req.Key = horus.UserByIdV(u.Id)
	return s.bare.User().Update(ctx, req)
}

func (s *UserServiceServer) Delete(ctx context.Context, req *horus.GetUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
