package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/horus/server"
)

type Server struct {
	Grpc         *grpc.Server
	GrpcListener net.Listener

	Rest         *http.Server
	RestListener net.Listener

	close_once sync.Once
}

func NewServer(h horus.Horus, conf *Config) (*Server, error) {
	errs := []error{}
	grpc_server, err := server.NewGrpcServer(h, conf.toGrpcServerConf())
	if err != nil {
		errs = append(errs, fmt.Errorf("create GRPC server: %w", err))
	}
	rest_server, err := server.NewRestServer(h, conf.toRestServerConf())
	if err != nil {
		errs = append(errs, fmt.Errorf("create REST server: %w", err))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	s := &Server{
		Grpc: grpc.NewServer(
			grpc.Creds(insecure.NewCredentials()),
			grpc.UnaryInterceptor(grpc_server.UnaryInterceptor),
		),
		Rest: &http.Server{
			Handler: httpLog(conf.Log.NewLogger(), rest_server),
		},
	}

	grpc_server.Register(s.Grpc)

	return s, nil
}

func (s *Server) Listen(conf *Config) error {
	errs := []error{}
	var err error = nil

	s.GrpcListener, err = net.Listen("tcp", conf.Grpc.Addr())
	if err != nil {
		errs = append(errs, err)
	}

	s.RestListener, err = net.Listen("tcp", conf.Rest.Addr())
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (s *Server) Serve(on_close func(err error)) {
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		err := s.Grpc.Serve(s.GrpcListener)
		if errors.Is(err, grpc.ErrServerStopped) {
			err = nil
		}
		on_close(err)
	}()
	go func() {
		defer wg.Done()
		err := s.Rest.Serve(s.RestListener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		on_close(err)
	}()
}

func (s *Server) Close(ctx context.Context) error {
	var err error
	s.close_once.Do(func() {
		err = s.close(ctx)
	})

	return err
}

func (s *Server) close(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()

	go func() {
		defer wg.Done()
		s.Grpc.GracefulStop()
	}()
	var err error
	go func() {
		defer wg.Done()
		err = s.Rest.Shutdown(ctx)
	}()

	return err
}
