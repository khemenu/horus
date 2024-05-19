package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/service/frame"
)

type UserService struct {
	horus.UnimplementedUserServiceServer
	*base
}

func (s *UserService) Create(ctx context.Context, req *horus.CreateUserRequest) (*horus.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

func (s *UserService) Get(ctx context.Context, req *horus.GetUserRequest) (*horus.User, error) {
	f := frame.Must(ctx)
	return s.bare.User().Get(ctx, &horus.GetUserRequest{
		Id:   f.Actor.ID[:],
		View: req.View,
	})
}

func (s *UserService) Update(ctx context.Context, req *horus.UpdateUserRequest) (*horus.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *UserService) Delete(ctx context.Context, req *horus.DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
