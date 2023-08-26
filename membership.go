package horus

import (
	"context"
	"time"
)

type Membership struct {
	TeamId   TeamId
	MemberId MemberId
	Role     RoleTeam

	CreatedAt time.Time
}

type MembershipInit struct {
	TeamId   TeamId
	MemberId MemberId
	Role     RoleTeam
}

type MembershipStore interface {
	New(ctx context.Context, init MembershipInit) (*Membership, error)
	GetById(ctx context.Context, team_id TeamId, member_id MemberId) (*Membership, error)
}
