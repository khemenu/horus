package env

import (
	"context"
	"errors"
	"fmt"
	"os"

	"khepri.dev/horus"
	"khepri.dev/horus/cmd/hr/reporter"
	"khepri.dev/horus/ent"
)

type Env struct {
	reporter reporter.Reporter

	// Make it private and make the methods to set it.
	TokenValue string
	ActorId    string

	db     *ent.Client
	server *buff_server
	client horus.Client
}

func (e *Env) Report(raw any, text string) error {
	if e.reporter == nil {
		e.reporter = &reporter.TextReporter{Out: os.Stdout}
	}

	if reporter.IsStructured(e.reporter) {
		return e.reporter.Report(raw)
	} else {
		return e.reporter.Report(text)
	}
}

func (e *Env) IsBareServer() bool {
	return e.TokenValue == "" && e.ActorId == ""
}

func (e *Env) ToBeBareServer() error {
	if !e.IsBareServer() {
		return fmt.Errorf("this operation requires the server to be a bare server")
	}

	return nil
}

func (e *Env) NotToBeBareServe() error {
	if e.IsBareServer() {
		return fmt.Errorf("this operation cannot be run with the bare server; please provide a token or an actor")
	}

	return nil
}

func (e *Env) CleanUp(ctx context.Context) error {
	errs := []error{}
	if e.server != nil {
		e.server.grpc_server.GracefulStop()
	}
	if e.db != nil {
		if err := e.db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close DB: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

type ctxKey struct{}

func From(ctx context.Context) *Env {
	v, ok := ctx.Value(ctxKey{}).(*Env)
	if !ok {
		panic("no env available")
	}

	return v
}

func Into(ctx context.Context, v *Env) context.Context {
	return context.WithValue(ctx, ctxKey{}, v)
}
