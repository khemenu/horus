package env

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/conf"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/log"
	"khepri.dev/horus/server"
	"khepri.dev/horus/server/bare"
	"khepri.dev/horus/server/frame"
)

func (e *Env) Connect(ctx context.Context) (horus.Client, error) {
	if e.client != nil {
		return e.client, nil
	}

	c := conf.From(ctx).Hr.Connect
	l := log.From(ctx)

	var (
		conn grpc.ClientConnInterface
		err  error
	)
	switch c.With {
	case "db":
		e.db, err = c.Db.Open()
		if err != nil {
			return nil, fmt.Errorf("open Ent client: %w", err)
		}
		if regexp.MustCompile(`^file\:.+\?mode\=memory`).MatchString(c.Db.Source) {
			l.Info("initialize mem DB")
			if err := e.db.Schema.Create(ctx); err != nil {
				return nil, fmt.Errorf("create DB schema: %w", err)
			}
		}

		var (
			horus_server     horus.Server
			grpc_server_opts []grpc.ServerOption

			covered_server = server.NewServer(e.db)
		)
		if e.IsBareServer() {
			l.Warn("bare server is enabled since no token or actor is provided")
			horus_server = &bareServer{
				covered: covered_server,
				Store:   bare.NewStore(e.db),
			}
		} else {
			f := &frame.Frame{}
			if e.TokenValue != "" {
				_, err := covered_server.Auth().TokenSignIn(frame.WithContext(ctx, f), &horus.TokenSignInRequest{
					Token: e.TokenValue,
				})
				if err != nil {
					return nil, fmt.Errorf("token sign in: %w", err)
				}
			} else if id, err := uuid.Parse(e.ActorId); err == nil {
				u, err := e.db.User.Get(ctx, id)
				if err != nil {
					return nil, fmt.Errorf("get actor by ID: %w", err)
				}
				f.Actor = u
			} else {
				u, err := e.db.User.Query().Where(user.AliasEQ(e.ActorId)).Only(ctx)
				if err != nil {
					return nil, fmt.Errorf("get actor by alias: %w", err)
				}
				f.Actor = u
			}

			horus_server = covered_server
			grpc_server_opts = []grpc.ServerOption{
				grpc.UnaryInterceptor(newFrameInjector(f)),
			}
		}

		grpc_server := grpc.NewServer(grpc_server_opts...)
		horus.GrpcRegister(grpc_server, horus_server)

		buff_server := &buff_server{
			listener:    bufconn.Listen(1024 * 1024),
			grpc_server: grpc_server,
		}
		buff_server.wg.Add(1)
		go func() {
			defer buff_server.wg.Done()
			buff_server.err = buff_server.grpc_server.Serve(buff_server.listener)
		}()

		conn, err = grpc.NewClient("passthrough://bufnet",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
				return buff_server.listener.DialContext(ctx)
			}),
		)

	case "horus":
		target := fmt.Sprintf("%s://%s", c.Horus.Schema, c.Horus.Addr)
		conn, err = grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	default:
		return nil, fmt.Errorf("unknown connection type")
	}
	if err != nil {
		return nil, fmt.Errorf("new Horus client: %w", err)
	}

	e.client = horus.NewClient(conn)
	return e.client, nil
}

func newFrameInjector(f *frame.Frame) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = frame.WithContext(ctx, f)
		return handler(ctx, req)
	}
}

type bareServer struct {
	covered horus.Server
	horus.Store
}

func (s *bareServer) Auth() horus.AuthServiceServer {
	return s.covered.Auth()
}

type buff_server struct {
	listener    *bufconn.Listener
	grpc_server *grpc.Server

	wg  sync.WaitGroup
	err error
}
