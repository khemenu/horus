package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"khepri.dev/horus"
	"khepri.dev/horus/log"
	"khepri.dev/horus/store/ent"
	"khepri.dev/horus/store/ent/token"
)

func _token(token *ent.Token) *horus.Token {
	return &horus.Token{
		Value:     token.ID,
		OwnerId:   token.OwnerID,
		Name:      token.Name,
		Type:      horus.TokenType(token.Type),
		CreatedAt: token.CreatedAt,
		ExpiredAt: token.ExpiredAt,
	}
}

type tokenStore struct {
	client    *ent.Client
	token_gen horus.Generator
}

type TokenStoreOption func(s *tokenStore) error

func WithTokenGenerator(generator horus.Generator) TokenStoreOption {
	return func(s *tokenStore) error {
		s.token_gen = generator
		return nil
	}
}

func NewTokenStore(client *ent.Client, opts ...TokenStoreOption) (horus.TokenStore, error) {
	s := &tokenStore{
		client:    client,
		token_gen: horus.DefaultOpaqueTokenGenerator,
	}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *tokenStore) Issue(ctx context.Context, init horus.TokenInit) (*horus.Token, error) {
	const MaxRetry = 3

	var (
		res *ent.Token
		err error
	)
	for i := 0; i < MaxRetry; i++ {
		var opaque string
		opaque, err = s.token_gen.New()
		if err != nil {
			return nil, fmt.Errorf("generate opaque token: %w", err)
		}

		now := time.Now().UTC()
		res, err = s.client.Token.Create().
			SetID(opaque).
			SetOwnerID(init.OwnerId).
			SetType(string(init.Type)).
			SetCreatedAt(now).
			SetExpiredAt(now.Add(init.Duration)).
			Save(ctx)

		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	log.FromCtx(ctx).Info("new token", "value", res.ID[0:8])
	return _token(res), nil
}

func (s *tokenStore) GetByValue(ctx context.Context, value string, token_type horus.TokenType) (*horus.Token, error) {
	res, err := s.client.Token.Query().
		Where(token.And(
			token.ID(value),
			token.ExpiredAtGT(time.Now().UTC()),
			token.Type(string(token_type)),
		)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, horus.ErrNotExist
		}

		return nil, fmt.Errorf("query: %w", err)
	}

	return _token(res), nil
}

func (s *tokenStore) Revoke(ctx context.Context, value string) error {
	now := time.Now()

	_, err := s.client.Token.UpdateOneID(value).
		Where(token.ExpiredAtGT(now)).
		SetExpiredAt(now).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil
		}

		return fmt.Errorf("query: %w", err)
	}

	return nil
}

func (s *tokenStore) RevokeAll(ctx context.Context, owner_id uuid.UUID) error {
	now := time.Now()

	_, err := s.client.Token.Update().
		Where(token.And(
			token.OwnerID(owner_id),
			token.ExpiredAtGTE(time.Now()),
		)).
		SetExpiredAt(now).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}
