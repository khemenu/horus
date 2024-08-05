package server

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/internal/entutils"
	"khepri.dev/horus/role"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

type InvitationServiceServer struct {
	horus.UnimplementedInvitationServiceServer
	*base
}

func (s *InvitationServiceServer) Create(ctx context.Context, req *horus.CreateInvitationRequest) (*horus.Invitation, error) {
	f := frame.Must(ctx)
	if req.GetType() != "internal" {
		return nil, status.Errorf(codes.Unimplemented, "type other than internal not implemented")
	}

	q := s.db.Account.Query().
		Where(account.HasOwnerWith(user.IDEQ(f.Actor.ID)))
	if p, err := bare.GetSiloSpecifier(req.GetSilo()); err != nil {
		return nil, err
	} else {
		q.Where(account.HasSiloWith(p))
	}

	actor_acct, err := q.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "account not found")
		}

		return nil, fmt.Errorf("get account: %w", err)
	}
	if actor_acct.Role != role.Owner {
		return nil, status.Error(codes.PermissionDenied, "the actor is not the owner of the silo")
	}

	invitee_id := req.GetInvitee()
	invitee_uuid, err := uuid.Parse(invitee_id)
	if err != nil {
		if invitee_id == "" {
			return nil, status.Errorf(codes.InvalidArgument, "invitee ID not provided")
		}

		invitee_uuid, err = s.db.User.Query().
			Where(user.AliasEQ(invitee_id)).
			OnlyID(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, status.Error(codes.NotFound, "invitee not found")
			}

			return nil, fmt.Errorf("get invitee: %w", err)
		}
	}

	ts_expired := req.GetDateExpired()
	if ts_expired == nil {
		ts_expired = timestamppb.New(time.Now().Add(7 * 24 * time.Hour))
	}

	return s.bare.Invitation().Create(ctx, &horus.CreateInvitationRequest{
		Invitee: invitee_uuid.String(),
		Type:    horus.InvitationTypeInternal,
		Silo:    req.GetSilo(),
		Inviter: horus.AccountById(actor_acct.ID),

		DateExpired: ts_expired,
	})
}

func (s *InvitationServiceServer) Get(ctx context.Context, req *horus.GetInvitationRequest) (*horus.Invitation, error) {
	f := frame.Must(ctx)

	v, err := s.bare.Invitation().Get(ctx, req)
	if err != nil {
		return nil, err
	}
	if v.Type == horus.InvitationTypeInternal && v.Invitee == f.Actor.ID.String() {
		return v, nil
	}

	actor_acct, err := s.db.Account.Query().
		Where(
			account.HasOwnerWith(user.IDEQ(f.Actor.ID)),
			account.HasSiloWith(silo.IDEQ(uuid.UUID(v.Silo.Id))),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.PermissionDenied, "account not found")
		}

		return nil, fmt.Errorf("get account: %w", err)
	}
	if actor_acct.Role != role.Owner {
		return nil, status.Error(codes.PermissionDenied, "only the silo owner can see the invitation")
	}

	return v, nil
}

func (s *InvitationServiceServer) Update(ctx context.Context, req *horus.UpdateInvitationRequest) (*horus.Invitation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *InvitationServiceServer) Accept(ctx context.Context, req *horus.AcceptInvitationRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)

	v, err := s.Get(ctx, horus.InvitationByIdV(req.GetId()))
	if err != nil {
		return nil, err
	}
	if v.DateExpired != nil && v.DateExpired.AsTime().Before(time.Now()) {
		return nil, status.Error(codes.PermissionDenied, "expired invitation")
	}
	if v.DateCanceled != nil {
		return nil, status.Error(codes.FailedPrecondition, "canceled invitation")
	}
	if v.DateDeclined != nil || v.DateAccepted != nil {
		return nil, status.Error(codes.FailedPrecondition, "invitation already settled")
	}
	if v.Type != horus.InvitationTypeInternal {
		return nil, status.Error(codes.Unimplemented, "type not implemented")
	}
	if v.Invitee != uuid.UUID(f.Actor.ID).String() {
		return nil, status.Error(codes.PermissionDenied, "only the invitee can accept the invitation")
	}

	return nil, entutils.WithTx(ctx, s.db, func(tx *ent.Tx) error {
		_, err := tx.Account.Create().
			SetRole(role.Owner).
			SetOwnerID(f.Actor.ID).
			SetSiloID(uuid.UUID(v.Silo.Id)).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("create an account: %w", err)
		}

		_, err = tx.Invitation.UpdateOneID(uuid.UUID(v.Id)).
			SetDateAccepted(time.Now()).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("update the invitation: %w", err)
		}

		return nil
	})
}

func (s *InvitationServiceServer) Delete(ctx context.Context, req *horus.GetInvitationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
