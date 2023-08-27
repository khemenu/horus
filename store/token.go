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

func fromEntToken(token *ent.Token) *horus.Token {
	return &horus.Token{
		Value:     token.ID,
		OwnerId:   horus.UserId(token.OwnerID),
		Name:      token.Name,
		Type:      horus.TokenType(token.Type),
		CreatedAt: token.CreatedAt,
		ExpiredAt: token.ExpiredAt,
	}
}

type tokenStore struct {
	*stores
}

func (s *tokenStore) Issue(ctx context.Context, init horus.TokenInit) (*horus.Token, error) {
	const MaxRetry = 3

	var (
		res *ent.Token
		err error
	)
	for i := 0; i < MaxRetry; i++ {
		var opaque string
		opaque, err = s.conf.TokenGenerator.New()
		if err != nil {
			return nil, fmt.Errorf("generate opaque token: %w", err)
		}

		now := time.Now().UTC()
		res, err = s.client.Token.Create().
			SetID(opaque).
			SetOwnerID(uuid.UUID(init.OwnerId)).
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
	return fromEntToken(res), nil
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

	return fromEntToken(res), nil
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

func (s *tokenStore) RevokeAll(ctx context.Context, owner_id horus.UserId) error {
	now := time.Now()

	_, err := s.client.Token.Update().
		Where(token.And(
			token.OwnerID(uuid.UUID(owner_id)),
			token.ExpiredAtGTE(time.Now()),
		)).
		SetExpiredAt(now).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}
