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

type MemberStore interface {
	New(ctx context.Context, init MemberInit) (*Member, error)
	GetById(ctx context.Context, member_id MemberId) (*Member, error)
	GetByUserIdFromOrg(ctx context.Context, org_id OrgId, user_id UserId) (*Member, error)
	GetAllByOrgId(ctx context.Context, org_id OrgId) ([]*Member, error)
	AddIdentity(ctx context.Context, member_id MemberId, identity_value IdentityValue) error
	UpdateById(ctx context.Context, member *Member) (*Member, error)
	DeleteById(ctx context.Context, member_id MemberId) error
	DeleteByUserIdFromOrg(ctx context.Context, org_id OrgId, user_id UserId) error
}
