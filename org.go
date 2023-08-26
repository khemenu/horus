package horus

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrgId uuid.UUID

func (i OrgId) String() string {
	return (uuid.UUID)(i).String()
}

type Org struct {
	Id   OrgId
	Name string

	CreatedAt time.Time
}

type OrgInit struct {
	OwnerId UserId
	Name    string
}

type OrgStore interface {
	New(ctx context.Context, init OrgInit) (*Org, error)
	GetById(ctx context.Context, org_id OrgId) (*Org, error)
	GetAllByUserId(ctx context.Context, user_id UserId) ([]*Org, error)
	UpdateById(ctx context.Context, org *Org) (*Org, error)
}
