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
	"khepri.dev/horus/server/frame"
)

type InvitationServiceServer struct {
	horus.UnimplementedInvitationServiceServer
	*base
}

func (s *InvitationServiceServer) Create(ctx context.Context, req *horus.CreateInvitationRequest) (*horus.Invitation, error) {
	f := frame.Must(ctx)
	if req.GetInvitation().GetType() != "internal" {
		return nil, status.Errorf(codes.Unimplemented, "type other than internal not implemented")
	}

	target_silo := req.GetInvitation().GetSilo()
	silo_uuid, err := uuid.FromBytes(target_silo.GetId())
	if err != nil && target_silo.GetId() != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid silo ID")
	}
	if silo_uuid == uuid.Nil {
		silo_alias := target_silo.GetAlias()
		if silo_alias == "" {
			return nil, status.Errorf(codes.InvalidArgument, "silo ID not provided")
		}

		silo_uuid, err = s.db.Silo.Query().
			Where(silo.AliasEQ(silo_alias)).
			OnlyID(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, status.Error(codes.NotFound, "silo not found")
			}

			return nil, fmt.Errorf("get silo: %w", err)
		}
	}

	actor_acct, err := s.db.Account.Query().
		Where(
			account.HasOwnerWith(user.IDEQ(f.Actor.ID)),
			account.HasSiloWith(silo.IDEQ(silo_uuid)),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "account not found")
		}

		return nil, fmt.Errorf("get account: %w", err)
	}
	if actor_acct.Role != account.RoleOWNER {
		return nil, status.Error(codes.PermissionDenied, "the actor is not the owner of the silo")
	}

	invitee_id := req.GetInvitation().GetInvitee()
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

	ts_expired := req.GetInvitation().GetDateExpired()
	if ts_expired == nil {
		ts_expired = timestamppb.New(time.Now().Add(7 * 24 * time.Hour))
	}

	return s.bare.Invitation().Create(ctx, &horus.CreateInvitationRequest{Invitation: &horus.Invitation{
		Invitee: invitee_uuid.String(),
		Type:    horus.InvitationTypeInternal,
		Silo:    &horus.Silo{Id: silo_uuid[:]},
		Inviter: &horus.Account{Id: actor_acct.ID[:]},

		DateExpired: ts_expired,
	}})
}

func (s *InvitationServiceServer) Get(ctx context.Context, req *horus.GetInvitationRequest) (*horus.Invitation, error) {
	f := frame.Must(ctx)

	req.View = horus.GetInvitationRequest_WITH_EDGE_IDS
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
	if actor_acct.Role != account.RoleOWNER {
		return nil, status.Error(codes.PermissionDenied, "only the silo owner can see the invitation")
	}

	return v, nil
}

func (s *InvitationServiceServer) Update(ctx context.Context, req *horus.UpdateInvitationRequest) (*horus.Invitation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (s *InvitationServiceServer) Accept(ctx context.Context, req *horus.AcceptInvitationRequest) (*emptypb.Empty, error) {
	f := frame.Must(ctx)

	v, err := s.Get(ctx, &horus.GetInvitationRequest{Id: req.GetId()})
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
			SetRole(account.RoleMEMBER).
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

func (s *InvitationServiceServer) Delete(ctx context.Context, req *horus.DeleteInvitationRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
