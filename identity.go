package horus

import (
	"context"
	"time"
)

type IdentityKind string

const (
	IdentityMail IdentityKind = "mail"
)

type Identity struct {
	OwnerId UserId
	Kind    IdentityKind
	Value   string
	Name    string

	VerifiedBy Verifier

	CreatedAt time.Time
}

type IdentityInit struct {
	OwnerId UserId
	Kind    IdentityKind
	Value   string
	Name    string

	VerifiedBy Verifier
}

type IdentityStore interface {
	New(ctx context.Context, input *IdentityInit) (*Identity, error)
	GetByValue(ctx context.Context, value string) (*Identity, error)
	GetAllByOwner(ctx context.Context, owner_id UserId) (map[string]*Identity, error)
	Update(ctx context.Context, input *Identity) (*Identity, error)
}
