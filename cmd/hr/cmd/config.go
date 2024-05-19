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
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DbConfig struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

type ClientConfig struct {
	Db ClientDbConfig `yaml:"db"`

	Svc horus.Service
}

type ClientDbConfig struct {
	Enabled bool   `yaml:"enabled"`
	UseBare bool   `yaml:"use_bare"`
	Driver  string `yaml:"driver"`
	Source  string `yaml:"source"`

	WithInit bool
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

	fx.Default(&c.Grpc.Host, "localhost")
	fx.Default(&c.Grpc.Port, 35122)
	fx.Default(&c.Log.Format, "text")
	fx.Default(&c.Client.Db.Driver, c.Db.Driver)
	fx.Default(&c.Client.Db.Source, c.Db.Source)

	if c.Grpc.Host == "0.0.0.0" {
		c.Grpc.Host = "localhost"
	}
	if c.Client.Db.Driver == "sqlite3" && strings.Contains(c.Client.Db.Source, "mode=memory") {
		c.Client.Db.WithInit = true
	}

	errs := []error{}
	if c.Client.Db.Enabled && (c.Client.Db.Driver == "" || c.Client.Db.Source == "") {
		errs = append(errs, fmt.Errorf(`".client.db.driver" or ".client.db.source" cannot be empty when ".client.db.enabled" is true`))
	}
	if !slices.Contains([]string{"text", "json"}, c.Log.Format) {
		errs = append(errs, fmt.Errorf(`log.format must be one of "text" or "json": %s`, c.Log.Format))
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return c, nil
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
