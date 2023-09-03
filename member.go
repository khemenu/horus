package horus

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MemberId uuid.UUID

func (i MemberId) String() string {
	return (uuid.UUID)(i).String()
}

type Member struct {
	Id     MemberId
	OrgId  OrgId
	UserId UserId
	Role   RoleOrg

	Name       string
	Identities map[string]*Identity

	CreatedAt time.Time
}

type MemberInit struct {
	OrgId  OrgId
	UserId UserId
	Role   RoleOrg
	Name   string
}

type MemberSortKeyword string

const (
	MemberSortByName        MemberSortKeyword = "name"
	MemberSortByCreatedDate MemberSortKeyword = "created_date"
)

type MemberSort struct {
	Keyword MemberSortKeyword
	Order   SortOrder
}

type MemberListFromOrgConfig struct {
	Offset int
	Limit  int
	Sorts  []MemberSort
}

type MemberStore interface {
	New(ctx context.Context, init MemberInit) (*Member, error)
	GetById(ctx context.Context, member_id MemberId) (*Member, error)
	GetByUserIdFromOrg(ctx context.Context, org_id OrgId, user_id UserId) (*Member, error)
	GetByUserIdFromTeam(ctx context.Context, team_id TeamId, user_id UserId) (*Member, error) // Resolved even if the user does not have a membership.
	GetAllByOrgId(ctx context.Context, org_id OrgId) ([]*Member, error)
	ListFromOrg(ctx context.Context, org_id OrgId, conf MemberListFromOrgConfig) ([]*Member, error)
	ListFromTeam(ctx context.Context, team_id TeamId, conf MemberListFromOrgConfig) ([]*Member, error)
	UpdateById(ctx context.Context, member *Member) (*Member, error)
	AddIdentity(ctx context.Context, member_id MemberId, identity_value IdentityValue) error
	AddIdentityByUserIdFromOrg(ctx context.Context, org_id OrgId, user_id UserId, identity_value IdentityValue) error
	RemoveIdentity(ctx context.Context, member_id MemberId, identity_value IdentityValue) error
	RemoveIdentityByUserIdFromOrg(ctx context.Context, org_id OrgId, user_id UserId, identity_value IdentityValue) error
	DeleteById(ctx context.Context, member_id MemberId) error
	DeleteByUserIdFromOrg(ctx context.Context, org_id OrgId, user_id UserId) error
}
