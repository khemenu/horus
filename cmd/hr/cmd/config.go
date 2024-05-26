package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"gopkg.in/yaml.v2"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/conf"
	"khepri.dev/horus/internal/fx"
)

type Config struct {
	path string

	Grpc GrpcConfig `yaml:"grpc"`

	Db     DbConfig     `yaml:"db"`
	Client ClientConfig `yaml:"client"`

	Log   conf.LogConfig `yaml:"log"`
	Debug DebugConfig    `yaml:"debug"`
}

type GrpcConfig struct {
	Network string `yaml:"network"`
	Address string `yaml:"address"`
}

type DbConfig struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

type ClientConfig struct {
	ConnectWith string `yaml:"connect_with"` // "target" | "db"

	Target ClientTargetConfig `yaml:"target"`
	Db     ClientDbConfig     `yaml:"db"`

	server *buff_server
	client horus.Client
}

type ClientTargetConfig struct {
	Schema  string `yaml:"schema"`
	Address string `yaml:"address"`
}

type ClientDbConfig struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`

	UseBare  bool `yaml:"use_bare"`
	WithInit bool `yaml:"with_init"`
}

type DebugConfig struct {
	Enabled   bool `yaml:"enabled"`
	Unsecured bool `yaml:"unsecured"`
}

func ReadConfig(path string) (*Config, error) {
	c := &Config{path: path}

	if f, err := os.Open(c.path); err != nil {
		return nil, fmt.Errorf("open config file at %s: %w", c.path, err)
	} else {
		err := yaml.NewDecoder(f).Decode(c)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("unmarshal config at %s: %w", c.path, err)
		}
	}

	return c, nil
}

func (c *Config) Evaluate() error {
	fx.Default(&c.Grpc.Network, "tcp")
	fx.Default(&c.Grpc.Address, "localhost:35122")
	fx.Default(&c.Client.ConnectWith, "target")
	fx.Default(&c.Client.Target.Schema, "dns")
	fx.Default(&c.Client.Target.Address, c.Grpc.Address)
	fx.Default(&c.Client.Db.Driver, c.Db.Driver)
	fx.Default(&c.Client.Db.Source, c.Db.Source)
	fx.Default(&c.Log.Enabled, fx.Addr(true))
	fx.Default(&c.Log.Format, "text")

	c.Client.Target.Address = strings.Replace(c.Client.Target.Address, "0.0.0.0", "localhost", 1)
	if c.Client.Db.Driver == "sqlite3" && strings.Contains(c.Client.Db.Source, "mode=memory") {
		c.Client.Db.WithInit = true
	}

	errs := []error{}
	if !slices.Contains([]string{"target", "db"}, c.Client.ConnectWith) {
		errs = append(errs, fmt.Errorf(`".client.connect_with" must be one of "target" or "db"`))
	}
	if !slices.Contains([]string{"text", "json"}, c.Log.Format) {
		errs = append(errs, fmt.Errorf(`log.format must be one of "text" or "json": %s`, c.Log.Format))
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}

type confCtxKey struct{}

func ConfFrom(ctx context.Context) *Config {
	l, ok := ctx.Value(confCtxKey{}).(*Config)
	if !ok {
		panic("no config available")
	}

	return l
}

func ConfInto(ctx context.Context, conf *Config) context.Context {
	return context.WithValue(ctx, confCtxKey{}, conf)
}

func (c *ClientConfig) CleanUp(ctx context.Context) error {
	if s := c.server; s != nil {
		s.grpc_server.GracefulStop()
		s.wg.Wait()
		if s.err != nil {
			return s.err
		}
	}

	return nil
}
