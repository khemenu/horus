package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/conf"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/log"
	"khepri.dev/horus/server"
	"khepri.dev/horus/server/frame"
	gw "khepri.dev/horus/server/gw"
)

var CmdServe = &cli.Command{
	Name:        "serve",
	Description: "Start Horus server",
	Action: func(ctx_ *cli.Context) error {
		ctx := ctx_.Context

		c := conf.From(ctx)
		l := log.From(ctx)
		if c.Debug.Enabled {
			l.Warn("debug mode is enabled")
		}

		var (
			db  *ent.Client
			err error
		)
		if c.Debug.Enabled && c.Debug.MemDb.Enabled {
			l.Warn("use mem DB")

			db, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
			if err != nil {
				return fmt.Errorf("create mem DB: %w", err)
			}
		} else {
			db, err = ent.Open(c.Db.Driver, c.Db.Source)
			if err != nil {
				return fmt.Errorf("create DB client: %w", err)
			}
		}
		if err := db.Schema.Create(ctx); err != nil {
			return fmt.Errorf("create DB schema: %w", err)
		}

		horus_server := server.NewServer(db)
		if c.Debug.Enabled && c.Debug.MemDb.Enabled {
			for _, u := range c.Debug.MemDb.Users {
				user, err := db.User.Create().SetAlias(u.Alias).Save(ctx)
				if err != nil {
					return fmt.Errorf("create user for mem DB: %w", err)
				}

				f := frame.Frame{Actor: user}
				ctx := frame.WithContext(ctx, &f)
				if _, err := horus_server.Token().Create(ctx, &horus.CreateTokenRequest{
					Value: u.Password,
					Type:  horus.TokenTypePassword,
				}); err != nil {
					return fmt.Errorf("set password for user %s: %w", u.Alias, err)
				}
			}
		}

		grpc_addr := fmt.Sprintf("%s:%d", c.Grpc.Host, c.Grpc.Port)
		http_addr := fmt.Sprintf("%s:%d", c.Http.Host, c.Http.Port)
		gw_addr := fmt.Sprintf("%s:%d", c.Grpc.Gateway.Host, c.Grpc.Gateway.Port)

		grpc_server := grpc.NewServer(
			grpc.Creds(insecure.NewCredentials()),
			grpc.ChainUnaryInterceptor(
				log.UnaryInterceptor(l, slog.LevelInfo),
				func() grpc.UnaryServerInterceptor {
					auth_interceptor := horus.AuthUnaryInterceptor(horus_server.Auth().TokenSignIn)
					svc_interceptor := server.UnaryInterceptor(horus_server, db)
					return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
						if strings.HasPrefix(info.FullMethod, "/khepri.horus.AuthService/") {
							return svc_interceptor(ctx, req, info, handler)
						}

						ctx = frame.WithContext(ctx, frame.New())
						return auth_interceptor(ctx, req, info, func(ctx context.Context, req any) (any, error) {
							return svc_interceptor(ctx, req, info, handler)
						})
					}
				}(),
			),
		)
		horus.GrpcRegister(grpc_server, horus_server)
		reflection.Register(grpc_server)

		http_server := &http.Server{}
		http_mux := http.NewServeMux()
		HandleAuth(http_mux, horus_server)
		HandleKubeWebhook(http_mux, horus_server)
		HandleRabbitMqHttpAuth(http_mux, horus_server)
		http_server.Handler = log.HttpLogger(l, slog.LevelInfo, http_mux)

		var gw_server *http.Server
		if c.Grpc.Gateway.Enabled {
			l.Info("gRPC gateway is enabled")
			opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

			gw_mux := runtime.NewServeMux()
			gw.RegisterConfServiceHandlerFromEndpoint(ctx, gw_mux, grpc_addr, opts)
			gw.RegisterUserServiceHandlerFromEndpoint(ctx, gw_mux, grpc_addr, opts)
			gw.RegisterIdentityServiceHandlerFromEndpoint(ctx, gw_mux, grpc_addr, opts)
			gw.RegisterAccountServiceHandlerFromEndpoint(ctx, gw_mux, grpc_addr, opts)
			gw.RegisterInvitationServiceHandlerFromEndpoint(ctx, gw_mux, grpc_addr, opts)
			gw.RegisterMembershipServiceHandlerFromEndpoint(ctx, gw_mux, grpc_addr, opts)
			gw.RegisterSiloServiceHandlerFromEndpoint(ctx, gw_mux, grpc_addr, opts)
			gw.RegisterTeamServiceHandlerFromEndpoint(ctx, gw_mux, grpc_addr, opts)
			gw.RegisterTokenServiceHandlerFromEndpoint(ctx, gw_mux, grpc_addr, opts)

			gw_server = &http.Server{Addr: gw_addr}
			gw_server.Handler = log.HttpLogger(l, slog.LevelInfo, gw_mux)
		}

		servers := ServerSet{}
		servers.Add("grpc", &GrpcServer{
			tcpListener: tcpListener{addr: grpc_addr},
			Server:      grpc_server,
		})
		l.Info("serve gRPC server", slog.String("addr", grpc_addr))

		servers.Add("http", &HttpServer{
			tcpListener: tcpListener{addr: http_addr},
			Server:      http_server,
		})
		l.Info("serve HTTP server", slog.String("addr", http_addr))

		if c.Grpc.Gateway.Enabled {
			servers.Add("grpc-gw", &HttpServer{
				tcpListener: tcpListener{addr: gw_addr},
				Server:      gw_server,
			})
			l.Info("serve gRPC HTTP gateway", slog.String("addr", gw_addr))
		}

		var serve_err error
		done := make(chan struct{}, 1)
		go func() {
			serve_err = servers.Serve()
			done <- struct{}{}
		}()

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

		select {
		case <-done:
			if serve_err != nil {
				return fmt.Errorf("serve: %w", serve_err)
			}
			return nil

		case sig := <-interrupt:
			l.Warn("interrupted", slog.String("signal", sig.String()))
			l.Warn("force shutdown after 1 minute; interrupt once more to shutdown now")
			go servers.Shutdown(ctx)
		}

		deadline := time.Now().Add(time.Minute)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		freq := false
	L:
		for {
			select {
			case <-done:
				l.Info("shutdown gracefully")
				return nil

			case sig := <-interrupt:
				l.Warn("interrupted", slog.String("signal", sig.String()))
				break L
			case <-ticker.C:
				r := time.Until(deadline)
				r = time.Duration((r+time.Second/2)/time.Second) * time.Second // floor
				if r < 15*time.Second && !freq {
					freq = true
					ticker.Reset(time.Second)
				}
				if r < 500*time.Millisecond {
					break L
				}

				l.Warn("tick", slog.Duration("remain", r))
			}
		}

		l.Warn("force shutdown")
		if err := servers.Close(); err != nil {
			l.Warn("force shutdown", slog.String("err", err.Error()))
		}

		return nil
	},
}
