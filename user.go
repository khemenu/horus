package horus

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id    uuid.UUID
	Alias string // Unique like ID but human readable.

	CreatedAt time.Time
}

type UserStore interface {
	New(ctx context.Context) (*User, error)
	GetById(ctx context.Context, id uuid.UUID) (*User, error)
	GetByAlias(ctx context.Context, alias string) (*User, error)
}
