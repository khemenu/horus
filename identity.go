package horus

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type IdentityKind string

const (
	IdentityEmail = IdentityKind("email")
)

type IdentityInit struct {
	Kind  IdentityKind
	Value string
	Name  string

	VerifiedBy Verifier
}

type Identity struct {
	IdentityInit
	OwnerId uuid.UUID

	CreatedAt time.Time
}

type IdentityStore interface {
	Create(ctx context.Context, input *Identity) (*Identity, error)
	GetByValue(ctx context.Context, value string) (*Identity, error)
	GetAllByOwner(ctx context.Context, owner_id uuid.UUID) (map[string]*Identity, error)
	Update(ctx context.Context, input *Identity) (*Identity, error)
}
