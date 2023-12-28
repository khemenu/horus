package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/exp/maps"
	"khepri.dev/horus"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/pb"
)

func toPbIdentityKind(v horus.IdentityKind) pb.IdentityKind {
	switch v {
	case horus.IdentityMail:
		return pb.IdentityKind_IDENTITY_KIND_MAIL
	}

	return pb.IdentityKind_IDENTITY_KIND_UNSPECIFIED
}

func toPbIdentity(v *horus.Identity) *pb.Identity {
	return &pb.Identity{
		OwnerId: v.OwnerId[:],
		Kind:    toPbIdentityKind(v.Kind),
		Value:   string(v.Value),

		Name:       v.Name,
		VerifiedBy: string(v.VerifiedBy),

		CreatedAt: v.CreatedAt.Format(time.RFC3339),
	}
}

func (s *grpcServer) ListIdentities(ctx context.Context, req *pb.ListIdentitiesReq) (*pb.ListIdentitiesRes, error) {
	user := s.mustUser(ctx)

	identities, err := s.Identities().GetAllByOwner(ctx, user.Id)
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("get identity details: %w", err))
	}

	return &pb.ListIdentitiesRes{
		Identities: fx.MapV(maps.Values(identities), toPbIdentity),
	}, nil
}

func (s *grpcServer) DeleteIdentity(ctx context.Context, req *pb.DeleteIdentityReq) (*pb.DeleteIdentityRes, error) {
	user := s.mustUser(ctx)

	identity, err := s.Identities().GetByValue(ctx, horus.IdentityValue(req.IdentityValue))
	if err != nil {
		if errors.Is(err, horus.ErrNotExist) {
			return &pb.DeleteIdentityRes{}, nil
		}

		return nil, grpcInternalErr(ctx, fmt.Errorf("get identity details: %w", err))
	}
	if identity.OwnerId != user.Id {
		return &pb.DeleteIdentityRes{}, nil
	}

	err = s.Identities().Delete(ctx, identity.Value)
	if err != nil {
		return nil, grpcInternalErr(ctx, fmt.Errorf("delete identity: %w", err))
	}

	return &pb.DeleteIdentityRes{}, nil
}
