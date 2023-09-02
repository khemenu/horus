package horus

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TeamId uuid.UUID

func (i TeamId) String() string {
	return (uuid.UUID)(i).String()
}

type Team struct {
	Id    TeamId
	OrgId OrgId
	Name  string

	CreatedAt time.Time
}

type TeamInit struct {
	OrgId   OrgId
	OwnerId MemberId
	Name    string
}

type TeamStore interface {
	New(ctx context.Context, init TeamInit) (*Team, error)
	GetById(ctx context.Context, team_id TeamId) (*Team, error)
	GetAllByOrgId(ctx context.Context, org_id OrgId) ([]*Team, error)
	UpdateById(ctx context.Context, team *Team) (*Team, error)
	DeleteByIdFromOrg(ctx context.Context, org_id OrgId, team_id TeamId) error
}
