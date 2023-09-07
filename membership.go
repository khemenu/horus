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
	GetByUserIdFromTeam(ctx context.Context, team_id TeamId, user_id UserId) (*Membership, error)
	GetAllByMemberId(ctx context.Context, member_id MemberId) ([]*Membership, error)
	UpdateById(ctx context.Context, membership *Membership) (*Membership, error)
	DeleteById(ctx context.Context, team_id TeamId, member_id MemberId) error
	DeleteByUserIdFromTeam(ctx context.Context, team_id TeamId, user_id UserId) error
}
