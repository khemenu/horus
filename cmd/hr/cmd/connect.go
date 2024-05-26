package cmd

import (
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/server"
	"khepri.dev/horus/server/bare"
)

func (c *ClientConfig) connect(ctx context.Context) (horus.Client, error) {
	if c.client != nil {
		return c.client, nil
	}

	var (
		conn grpc.ClientConnInterface
		err  error
	)
	switch c.ConnectWith {
	case "target":
		target := fmt.Sprintf("%s://%s", c.Target.Schema, c.Target.Address)
		conn, err = grpc.DialContext(ctx, target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

	case "db":
		var db *ent.Client
		db, err = ent.Open(c.Db.Driver, c.Db.Source)
		if err != nil {
			return nil, fmt.Errorf("open Ent client: %w", err)
		}
		if c.Db.WithInit {
			if err := db.Schema.Create(ctx); err != nil {
				return nil, fmt.Errorf("create DB schema: %w", err)
			}
		}

		horus_server := server.NewServer(db)
		if c.Db.UseBare {
			horus_server = &bare_server{
				db:    db,
				Store: bare.NewStore(db),
				cover: horus_server,
			}
		}

		buff_server := &buff_server{
			listener:    bufconn.Listen(1024 * 1024),
			grpc_server: grpc.NewServer(), // TODO: interceptor
		}

		horus.GrpcRegister(buff_server.grpc_server, horus_server)
		buff_server.wg.Add(1)
		go func() {
			defer buff_server.wg.Done()
			if err := buff_server.grpc_server.Serve(buff_server.listener); err != nil {
				buff_server.err = fmt.Errorf("gRPC serve: %w", err)
			}
		}()

		conn, err = grpc.DialContext(ctx, "",
			grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
				return buff_server.listener.DialContext(ctx)
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	default:
		panic("unknown connection type")
	}
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}

	c.client = horus.NewClient(conn)
	return c.client, nil
}

func (c *ClientConfig) mustConnect(ctx context.Context) horus.Client {
	return fx.Must(c.connect(ctx))
}

func (c *ClientConfig) mustBareServer(ctx context.Context) *bare_server {
	c.mustConnect(ctx)
	if c.server == nil {
		panic("client must be connected with DB for this operation")
	}

	server, ok := c.server.Server.(*bare_server)
	if !ok {
		panic("server must be bare server for this operation")
	}

	return server
}

type bare_server struct {
	db *ent.Client
	horus.Store
	cover horus.Server
}

func (s *bare_server) Auth() horus.AuthServiceServer {
	return s.cover.Auth()
}

type buff_server struct {
	horus.Server

	listener    *bufconn.Listener
	grpc_server *grpc.Server
	wg          sync.WaitGroup

	err error
}
