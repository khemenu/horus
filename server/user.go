package server

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
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
	f := frame.Must(ctx)

	u, err := s.bare.User().Get(ctx, req)
	if err != nil {
		return nil, err
	}

	// Is the target user an actor or a descendant of the actor?
	ok := false
	cursor := u
	for {
		ok = bytes.Equal(cursor.Id, f.Actor.ID[:])
		if ok {
			break
		}
		if cursor.Parent == nil {
			break
		}

		cursor, err = s.bare.User().Get(ctx, horus.UserByIdV(cursor.Parent.Id))
		if err != nil {
			return nil, err
		}
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	}
	if u.Parent == nil {
		return nil, status.Error(codes.FailedPrecondition, "root user cannot be deleted")
	}

	err = entutils.WithTx(ctx, s.db, func(tx *ent.Tx) error {
		// Ensure the user still exist.
		_, err := tx.User.Get(ctx, uuid.UUID(u.Id))
		if err != nil {
			return err
		}

		err = tx.User.Update().
			Where(user.HasParentWith(user.IDEQ(uuid.UUID(u.Id)))).
			SetParentID(uuid.UUID(u.Parent.Id)).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("updates parent of children: %w", err)
		}

		err = tx.User.DeleteOneID(uuid.UUID(u.Id)).Exec(ctx)
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
