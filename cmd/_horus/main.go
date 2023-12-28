package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	_ "github.com/mattn/go-sqlite3"
	"khepri.dev/horus/cmd/horus/server/app"
	"khepri.dev/horus/store"
	"khepri.dev/horus/store/ent"
)

func run(conf *Config) error {
	l := conf.Log.NewLogger()
	l.Info("read config", slog.String("path", conf.path))
	if conf.Debug.Enabled {
		l.Warn("debug mode is enabled")
	}
	if conf.Providers.IsEmpty() {
		l.Warn("no providers are configured")
	}

	var (
		client *ent.Client
		err    error
	)
	if conf.Debug.Enabled && conf.Debug.UseMemDb {
		l.Warn("use mem DB")
		client, err = store.NewSqliteMemClient()
		if err != nil {
			return fmt.Errorf("create mem DB client: %w", err)
		}
	} else {
		client, err = ent.Open(conf.Db.Driver, conf.Db.Source)
		if err != nil {
			return fmt.Errorf("create DB client: %w", err)
		}
	}

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	stores, err := store.NewStores(client, nil)
	if err != nil {
		return fmt.Errorf("create stores: %w", err)
	}

	horus, err := app.NewHorus(stores, conf.toHorusConf())
	if err != nil {
		return fmt.Errorf("create horus: %w", err)
	}

	server, err := NewServer(horus, conf)
	if err != nil {
		return fmt.Errorf("create server: %w", err)
	}
	l.Info("listen REST", "addr", conf.Rest.Addr())
	l.Info("listen GRPC", "addr", conf.Grpc.Addr())
	if err := server.Listen(conf); err != nil {
		return fmt.Errorf("server listen: %w", err)
	}

	on_exit := make(chan struct{})
	defer close(on_exit)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		shutdown_requested := false
		for {
			select {
			case <-on_exit:
				return

			case sig := <-signals:
				if shutdown_requested {
					l.Error("force the program to exit")
					os.Exit(1)
				}

				shutdown_requested = true
				l.Warn("shutting down the server", "signal", sig)
				server.Close(ctx)
			}
		}
	}()

	server.Serve(func(err error) {
		if err != nil {
			l.Error("unexpected close of server", "err", err)
		}

		go func() {
			if err := server.Close(ctx); err != nil {
				l.Error("error while server close", "err", err)
			}
		}()
	})

	return nil
}

func main() {
	conf, err := ParseArgs(os.Args)
	if err != nil {
		panic(err)
	}

	if err := run(conf); err != nil {
		conf.Log.NewLogger().Error("fatal", "err", err)
		os.Exit(1)
	}
}
