package entutils

import (
	"context"
	"errors"
	"fmt"

	"khepri.dev/horus/ent"
)

func WithTxV[T any](ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) (*T, error)) (*T, error) {
	tx, err := client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	res, err := fn(tx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Join(err, fmt.Errorf("%w: rolling back transaction: %v", err, rerr))
		}
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}
	return res, nil
}

func WithTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	_, err := WithTxV(ctx, client, func(tx *ent.Tx) (*int, error) {
		return nil, fn(tx)
	})
	return err
}
