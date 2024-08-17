package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"google.golang.org/grpc"
	"khepri.dev/horus/internal/fx"
)

type Server interface {
	Listen() (net.Listener, error)
	Serve(lis net.Listener) error
	Shutdown(ctx context.Context) error
	Close() error
}

type ServerState struct {
	Name string

	Server   Server
	Listener net.Listener
}

type ServerSet struct {
	wg sync.WaitGroup

	serving   sync.Mutex
	stopping  sync.Mutex
	closing   sync.Mutex
	close_err error
	is_closed bool

	servers []*ServerState
}

func (s *ServerSet) Err() error {
	s.closing.Lock()
	defer s.closing.Unlock()

	return s.close_err
}

func (s *ServerSet) Add(name string, svr Server) bool {
	s.wg.Add(1)

	for _, state := range s.servers {
		if state.Name == name {
			return false
		}
	}

	s.servers = append(s.servers, &ServerState{
		Name:   name,
		Server: svr,
	})

	return true
}

func (s *ServerSet) listen() error {
	errs := []error{}
	for _, state := range s.servers {
		l, err := state.Server.Listen()
		if err != nil {
			errs = append(errs, fmt.Errorf("listen %s: %w", state.Name, err))
		}

		state.Listener = l
	}
	if len(errs) == 0 {
		return nil
	}

	for _, state := range s.servers {
		if state.Listener == nil {
			continue
		}

		err := state.Listener.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("close listening \"%s\": %w", state.Name, err))
		}
	}
	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func (s *ServerSet) Serve() error {
	s.serving.Lock()
	defer s.serving.Unlock()

	if err := s.listen(); err != nil {
		return err
	}

	errs := []error{}
	for _, state := range s.servers {
		go func() {
			defer s.wg.Done()
			err := state.Server.Serve(state.Listener)
			if err != nil {
				errs = append(errs, fmt.Errorf("serve \"%s\": %w", state.Name, err))
			}

			s.Close()
		}()
	}

	s.wg.Wait()
	return errors.Join(errs...)
}

func (s *ServerSet) stop(ctx context.Context) {
	for _, state := range fx.Reversed(s.servers) {
		if err := state.Server.Shutdown(ctx); err != nil {
			s.Close()
			return
		}
	}
}

func (s *ServerSet) Shutdown(ctx context.Context) error {
	s.stopping.Lock()
	defer s.stopping.Unlock()

	s.stop(ctx)
	return s.close_err
}

func (s *ServerSet) close() error {
	errs := []error{}
	for _, state := range s.servers {
		err := state.Server.Close()
		if err != nil {
			errs = append(errs, fmt.Errorf("close server \"%s\": %w", state.Name, err))
		}
	}

	return errors.Join(errs...)
}

func (s *ServerSet) Close() error {
	s.closing.Lock()
	defer s.closing.Unlock()

	if !s.is_closed {
		s.is_closed = true
		s.close_err = s.close()
	}

	return s.close_err
}

type tcpListener struct {
	addr string
}

func (l *tcpListener) Listen() (net.Listener, error) {
	return net.Listen("tcp", l.addr)
}

type HttpServer struct {
	tcpListener
	*http.Server
}

func (s *HttpServer) Serve(lis net.Listener) error {
	err := s.Server.Serve(lis)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

type GrpcServer struct {
	tcpListener
	*grpc.Server
}

func (s *GrpcServer) Serve(lis net.Listener) error {
	err := s.Server.Serve(lis)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}

	return nil
}

func (s *GrpcServer) Shutdown(ctx context.Context) error {
	s.Server.GracefulStop()
	return nil
}

func (s *GrpcServer) Close() error {
	s.Server.Stop()
	return nil
}
