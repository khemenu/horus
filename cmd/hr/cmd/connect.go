package cmd

import (
	"context"
	"fmt"

	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/server"
	"khepri.dev/horus/server/bare"
)

func (c *ClientConfig) connect(ctx context.Context) (horus.Server, error) {
	if c.Svc != nil {
		return c.Svc, nil
	}

	db, err := ent.Open(c.Db.Driver, c.Db.Source)
	if err != nil {
		return nil, fmt.Errorf("open Ent client: %w", err)
	}

	if c.Db.WithInit {
		if err := db.Schema.Create(ctx); err != nil {
			return nil, fmt.Errorf("create DB schema: %w", err)
		}
	}

	svc := server.NewServer(db)
	if c.Db.UseBare {
		c.Svc = &bare_service{
			client: db,
			Store:  bare.NewStore(db),
			svc:    svc,
		}
	} else {
		c.Svc = svc
	}

	return c.Svc, nil
}

func (c *ClientConfig) mustConnect(ctx context.Context) horus.Server {
	return fx.Must(c.connect(ctx))
}

func (c *ClientConfig) mustBareConnect(ctx context.Context) *bare_service {
	s := c.mustConnect(ctx)

	bare_svc, ok := s.(*bare_service)
	if !ok {
		panic("service must be bare service for this operation")
	}

	return bare_svc
}

type bare_service struct {
	client *ent.Client
	horus.Store
	svc horus.Server
}

func (s *bare_service) Auth() horus.AuthServiceServer {
	return s.svc.Auth()
}
