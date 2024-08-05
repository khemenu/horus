package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/server/bare"
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
	if req.GetAlias() == "" {
		req.Alias = nil
	}
	req.Parent = horus.UserById(f.Actor.ID)
	return s.bare.User().Create(ctx, req)
}

func (s *UserServiceServer) Get(ctx context.Context, req *horus.GetUserRequest) (*horus.User, error) {
	f := frame.Must(ctx)
	fmt.Printf("req: %v\n", req)
	v, err := s.hasPermission(ctx, f.Actor, req)
	if err != nil {
		return nil, err
	}

	req = horus.UserById(v.ID)
	return s.bare.User().Get(ctx, req)
}

func (s *UserServiceServer) Update(ctx context.Context, req *horus.UpdateUserRequest) (*horus.User, error) {
	f := frame.Must(ctx)

	// TODO: Enable child transfer?
	if req.GetParent() != nil {
		req.Parent = nil
	}

	v, err := s.hasPermission(ctx, f.Actor, req.GetKey())
	if err != nil {
		return nil, err
	}

	req.Key = horus.UserById(v.ID)
	return s.bare.User().Update(ctx, req)
}

func (s *UserServiceServer) Delete(ctx context.Context, req *horus.GetUserRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)

	v, err := s.hasPermission(ctx, f.Actor, req)
	if err != nil {
		return nil, err
	}
	if v.ParentID == nil {
		return nil, status.Error(codes.FailedPrecondition, "root user cannot be deleted")
	}

	err = entutils.WithTx(ctx, s.db, func(tx *ent.Tx) error {
		// Ensure the user still exist.
		_, err := tx.User.Get(ctx, v.ID)
		if err != nil {
			return err
		}

		err = tx.User.Update().
			Where(user.HasParentWith(user.IDEQ(v.ID))).
			SetParentID(*v.ParentID).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("updates parent of children: %w", err)
		}

		err = tx.User.DeleteOneID(v.ID).Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, bare.ToStatus(err)
	}

	return &emptypb.Empty{}, nil
}
