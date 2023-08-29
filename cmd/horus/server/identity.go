package server

import (
	"context"
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
