package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/log"
	"khepri.dev/horus/service"
)

func Run(ctx context.Context, c *Config) error {
	l := c.Log.NewLogger()
	l.Info("config read", slog.String("path", c.path))
	if c.Debug.Enabled {
		l.Warn("debug mode is enabled")
	}

	ctx = log.WithCtx(ctx, l)
	var (
		db  *ent.Client
		err error
	)
	if c.Debug.Enabled && c.Debug.UseMemDb {
		l.Warn("use mem DB")

		db, err = ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
		if err != nil {
			return fmt.Errorf("create mem DB db: %w", err)
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
		grpc.UnaryInterceptor(service.GrpcUnaryInterceptor(svc, db)))
	http_server := &http.Server{
		Addr: http_addr,
	}
	shutdown := func() {
		grpc_server.GracefulStop()
		http_server.Shutdown(ctx)
	}

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

		service.GrpcRegisterStoreService(svc, grpc_server)
		reflection.Register(grpc_server)

		l.Info("serve gRPC", slog.String("addr", grpc_addr))
		err_grpc = grpc_server.Serve(grpc_listener)
	}()
	go func() {
		defer wg.Done()
		defer once.Do(shutdown)

		mux := http.NewServeMux()
		HandleAuth(mux, svc)
		http_server.Handler = mux

		l.Info("serve HTTP", slog.String("addr", http_addr))
		err_http = http_server.Serve(http_listener)
	}()

	errs := []error{}
	wg.Wait()
	if err_grpc != nil && !errors.Is(err, grpc.ErrServerStopped) {
		errs = append(errs, fmt.Errorf("unexpected gRPC server stop: %w", err_grpc))
	}
	if err_http != nil && !errors.Is(err, http.ErrServerClosed) {
		errs = append(errs, fmt.Errorf("unexpected HTTP server stop: %w", err_http))
	}

	return errors.Join(errs...)
}
