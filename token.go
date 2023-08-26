package horus

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenType string

var (
	RefreshToken TokenType = "refresh_token"
	AccessToken  TokenType = "access_token"
)

type TokenInit struct {
	OwnerId UserId

	Type     TokenType
	Name     string
	Duration time.Duration
}

type Token struct {
	Value   string
	OwnerId UserId

	Type TokenType
	Name string

	CreatedAt time.Time
	ExpiredAt time.Time
}

func (t *Token) Duration() time.Duration {
	return t.ExpiredAt.Sub(t.CreatedAt)
}

type TokenStore interface {
	Issue(ctx context.Context, init TokenInit) (*Token, error)
	GetByValue(ctx context.Context, value string, token_type TokenType) (*Token, error)
	Revoke(ctx context.Context, value string) error
	RevokeAll(ctx context.Context, owner_id uuid.UUID) error
}
