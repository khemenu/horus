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

	Name     string
	Contacts map[string]*Identity

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
	GetByUserIdFromOrg(ctx context.Context, org_id OrgId, user_id UserId) (*Member, error)
	GetAllByOrgId(ctx context.Context, org_id OrgId) ([]*Member, error)
}
