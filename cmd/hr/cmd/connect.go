package cmd

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/log"
	"khepri.dev/horus/server"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

func (c *ClientConfig) connect(ctx context.Context) (horus.Client, error) {
	if c.client != nil {
		return c.client, nil
	}

	l := log.From(ctx)
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
		c.db, err = ent.Open(c.Db.Driver, c.Db.Source)
		if err != nil {
			return nil, fmt.Errorf("open Ent client: %w", err)
		}
		if c.Db.WithInit {
			if err := c.db.Schema.Create(ctx); err != nil {
				return nil, fmt.Errorf("create DB schema: %w", err)
			}
		}

		var horus_server horus.Server
		if c.isBareServer() {
			l.Warn("bare server is enabled since no token or actor is provided")
			horus_server = &bare_server{
				covered: server.NewServer(c.db),
				Store:   bare.NewStore(c.db),
			}
		} else {
			horus_server = &covered_server{
				bare:   bare.NewStore(c.db),
				Server: server.NewServer(c.db),
			}
		}

		grpc_server := grpc.NewServer(grpc.ChainUnaryInterceptor(
			func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
				if vs := metadata.ValueFromIncomingContext(ctx, "actor-uuid"); len(vs) > 0 {
					id, err := uuid.Parse(vs[0])
					if err != nil {
						return nil, status.Errorf(codes.Internal, "invalid actor UUID")
					}
					user, err := c.db.User.Get(ctx, id)
					if err != nil {
						return nil, status.Errorf(codes.Internal, "get user for actor")
					}

					ctx = frame.WithContext(ctx, &frame.Frame{
						Actor: user,
					})
				}

				return handler(ctx, req)
			},
		))
		horus.GrpcRegister(grpc_server, horus_server)

		buff_server := &buff_server{
			listener:     bufconn.Listen(1024 * 1024),
			horus_server: horus_server,
			grpc_server:  grpc_server,
		}
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

func (c *ClientConfig) connectDbServer(ctx context.Context) (horus.Client, error) {
	if c.ConnectWith != "db" {
		return nil, fmt.Errorf(`".client.connect_with" must be "db" for this operation`)
	}

	return c.connect(ctx)
}

type bare_server struct {
	covered horus.Server
	horus.Store
}

func (s *bare_server) Auth() horus.AuthServiceServer {
	return s.covered.Auth()
}

type covered_server struct {
	bare horus.Store
	horus.Server
}

type buff_server struct {
	listener     *bufconn.Listener
	grpc_server  *grpc.Server
	horus_server horus.Server

	wg  sync.WaitGroup
	err error
}
