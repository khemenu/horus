package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"khepri.dev/horus"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/log"
	"khepri.dev/horus/service"
	"khepri.dev/horus/service/frame"
)

func Run(ctx context.Context, c *Config) error {
	l := c.Log.NewLogger()
	l.Info("config read", slog.String("path", c.path))
	if c.Debug.Enabled {
		l.Warn("debug mode is enabled")
	}

	ctx = log.Into(ctx, l)
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

	svc := service.NewService(db)
	if c.Debug.Enabled && c.Debug.MemDb.Enabled {
		for _, u := range c.Debug.MemDb.Users {
			user, err := db.User.Create().SetName(u.Name).Save(ctx)
			if err != nil {
				return fmt.Errorf("create user for mem DB: %w", err)
			}

			f := frame.Frame{Actor: user}
			ctx := frame.WithContext(ctx, &f)
			if _, err := svc.Token().Create(ctx, &horus.CreateTokenRequest{Token: &horus.Token{
				Value: u.Password,
				Type:  horus.TokenTypeBasic,
			}}); err != nil {
				return fmt.Errorf("set password for user %s: %w", u.Name, err)
			}
		}
	}

	grpc_addr := fmt.Sprintf("%s:%d", c.Grpc.Host, c.Grpc.Port)
	grpc_listener, err := net.Listen("tcp", grpc_addr)
	if err != nil {
		return fmt.Errorf("listen for gGRP: %w", err)
	}
	defer grpc_listener.Close()

	http_addr := fmt.Sprintf("%s:%d", c.Http.Host, c.Http.Port)
	http_listener, err := net.Listen("tcp", http_addr)
	if err != nil {
		return fmt.Errorf("listen for HTTP: %w", err)
	}
	defer http_listener.Close()

	grpc_server := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			log.UnaryInterceptor(l, slog.LevelInfo),
			func() grpc.UnaryServerInterceptor {
				auth_interceptor := horus.AuthUnaryInterceptor(svc.Auth().TokenSignIn)
				svc_interceptor := service.UnaryInterceptor(svc, db)
				return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
					if strings.HasPrefix(info.FullMethod, "/khepri.horus.AuthService/") {
						return svc_interceptor(ctx, req, info, handler)
					}

					return auth_interceptor(ctx, req, info, func(ctx context.Context, req any) (any, error) {
						return svc_interceptor(ctx, req, info, handler)
					})
				}
			}(),
		),
	)
	horus.GrpcRegister(svc, grpc_server)
	reflection.Register(grpc_server)

	http_server := &http.Server{
		Addr: http_addr,
	}
	http_mux := http.NewServeMux()
	HandleAuth(http_mux, svc)
	http_server.Handler = log.HttpLogger(l, slog.LevelInfo, http_mux)

	shutdown := func() {
		grpc_server.GracefulStop()
		http_server.Shutdown(ctx)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var (
		wg   sync.WaitGroup
		once sync.Once

		err_grpc error
		err_http error
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		defer once.Do(shutdown)

		l.Info("serve gRPC", slog.String("addr", grpc_addr))
		err_grpc = grpc_server.Serve(grpc_listener)
	}()
	go func() {
		defer wg.Done()
		defer once.Do(shutdown)

		l.Info("serve HTTP", slog.String("addr", http_addr))
		err_http = http_server.Serve(http_listener)
	}()

	graceful := make(chan struct{}, 1)
	go func() {
		select {
		case sig := <-interrupt:
			l.Warn("interrupted", slog.String("signal", sig.String()))
		case <-graceful:
			return
		}
		once.Do(shutdown)

		l.Warn("force shutdown after 1 minute. interrupt once more to shutdown now.")
		deadline := time.Now().Add(time.Minute)
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		freq := false
	L:
		for {
			select {
			case <-graceful:
				return

			case <-interrupt:
				break L
			case <-ticker.C:
				r := time.Until(deadline)
				r = time.Duration((r+time.Second/2)/time.Second) * time.Second
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

		l.Error("force shutdown")
		os.Exit(1)
	}()

	errs := []error{}
	wg.Wait()
	if err_grpc != nil && !errors.Is(err_grpc, grpc.ErrServerStopped) {
		errs = append(errs, fmt.Errorf("unexpected gRPC server stop: %w", err_grpc))
	}
	if err_http != nil && !errors.Is(err_http, http.ErrServerClosed) {
		errs = append(errs, fmt.Errorf("unexpected HTTP server stop: %w", err_http))
	}

	close(graceful)
	return errors.Join(errs...)
}
