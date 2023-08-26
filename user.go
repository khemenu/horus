package horus

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserId uuid.UUID

func (i UserId) String() string {
	return (uuid.UUID)(i).String()
}

type User struct {
	Id    UserId
	Alias string // Unique like ID but human readable.

	CreatedAt time.Time
}

type UserStore interface {
	New(ctx context.Context) (*User, error)
	GetById(ctx context.Context, user_id UserId) (*User, error)
	GetByAlias(ctx context.Context, alias string) (*User, error)
}
