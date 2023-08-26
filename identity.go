package horus

import (
	"context"
	"time"
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
	OwnerId UserId

	CreatedAt time.Time
}

type IdentityStore interface {
	Create(ctx context.Context, input *Identity) (*Identity, error)
	GetByValue(ctx context.Context, value string) (*Identity, error)
	GetAllByOwner(ctx context.Context, owner_id UserId) (map[string]*Identity, error)
	Update(ctx context.Context, input *Identity) (*Identity, error)
}
