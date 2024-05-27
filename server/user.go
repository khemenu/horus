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
	return s.bare.User().Create(ctx, &horus.CreateUserRequest{User: &horus.User{
		Alias: req.GetUser().GetAlias(),
		Parent: &horus.User{
			Id: f.Actor.ID[:],
		},
	}})
}

func (s *UserServiceServer) Get(ctx context.Context, req *horus.GetUserRequest) (*horus.User, error) {
	f := frame.Must(ctx)
	id := req.GetId()
	if id == nil {
		id = f.Actor.ID[:]
	}
	v, err := s.bare.User().Get(ctx, &horus.GetUserRequest{
		Id:   id,
		View: horus.GetUserRequest_WITH_EDGE_IDS,
	})
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
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *UserServiceServer) Delete(ctx context.Context, req *horus.DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
