package conf

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/internal/fx"
	"khepri.dev/horus/log"
)

type Config struct {
	Path string

	Grpc GrpcConfig `yaml:"grpc"`
	Http HttpConfig `yaml:"http"`

	Db DbConfig `yaml:"db"`
	Hr HrConfig `yaml:"hr"`

	Log   LogConfig   `yaml:"log"`
	Debug DebugConfig `yaml:"debug"`
}

type GrpcConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	Gateway GrpcGwConfig `yaml:"gateway"`
}

type GrpcGwConfig struct {
	Enabled    bool `yaml:"enabled"`
	HttpConfig `yaml:",inline"`
}

type HttpConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DbConfig struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

func (c *DbConfig) Open() (*ent.Client, error) {
	return ent.Open(c.Driver, c.Source)
}

type HrConfig struct {
	Connect HrConnectConfig `yaml:"connect"`
	// Reporter ReporterConfig  `yaml:"reporter"`
}

type HrConnectConfig struct {
	With string `yaml:"with"` // "db" | "horus"

	Db    DbConfig             `yaml:"db"`
	Horus HrConnectHorusConfig `yaml:"horus"`
}

type HrConnectHorusConfig struct {
	Schema string `yaml:"schema"` // "http" | "https"
	Addr   string `yaml:"addr"`
}

type DebugConfig struct {
	Enabled   bool        `yaml:"enabled"`
	Unsecured bool        `yaml:"unsecured"`
	MemDb     MemDbConfig `yaml:"mem_db"`
}

type MemDbConfig struct {
	Enabled bool `yaml:"enabled"`
	Users   []struct {
		Alias    string `yaml:"alias"`
		Password string `yaml:"password"`
	} `yaml:"users"`
}

func (c *Config) Evaluate() error {
	fx.Default(&c.Grpc.Host, "0.0.0.0")
	fx.Default(&c.Grpc.Port, 35122)
	fx.Default(&c.Http.Host, "0.0.0.0")
	fx.Default(&c.Http.Port, 35123)

	grpc_addr := fmt.Sprintf("%s:%d", c.Grpc.Host, c.Grpc.Port)
	grpc_addr = strings.Replace(grpc_addr, "0.0.0.0", "localhost", 1)

	fx.Default(&c.Db.Driver, "sqlite3")
	fx.Default(&c.Db.Source, "file:horus.db?cache=shared&_fk=1")
	fx.Default(&c.Hr.Connect.With, "db")
	fx.Default(&c.Hr.Connect.Db.Driver, c.Db.Driver)
	fx.Default(&c.Hr.Connect.Db.Source, c.Db.Source)
	fx.Default(&c.Hr.Connect.Horus.Schema, "http")
	fx.Default(&c.Hr.Connect.Horus.Addr, grpc_addr)

	fx.Default(&c.Log.Enabled, fx.Addr(true))
	fx.Default(&c.Log.Format, "text")
	fx.Default(&c.Log.Level, slog.LevelInfo)

	errs := []error{}
	if !slices.Contains([]string{"db", "horus"}, c.Hr.Connect.With) {
		errs = append(errs, fmt.Errorf(`".hr.connect.with" must be one of "db" or "target": %s`, c.Hr.Connect.With))
	}
	if !slices.Contains([]string{"text", "json"}, c.Log.Format) {
		errs = append(errs, fmt.Errorf(`.log.format must be one of "text" or "json": %s`, c.Log.Format))
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}

type ctxKey struct{}

func From(ctx context.Context) *Config {
	v, ok := ctx.Value(ctxKey{}).(*Config)
	if !ok {
		panic("no config available")
	}

	return v
}

func Into(ctx context.Context, v *Config) context.Context {
	return context.WithValue(ctx, ctxKey{}, v)
}

func FromFile(path string) (*Config, error) {
	c := &Config{Path: path}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return c, nil
}

func InitCmd(ctx *cli.Context, transform func(c *Config) error) (*Config, error) {
	p := ctx.String("conf")
	c, err := FromFile(p)
	if err != nil {
		return nil, fmt.Errorf("read config at %s: %w", p, err)
	}

	if err := transform(c); err != nil {
		return nil, err
	}
	if err := c.Evaluate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	l := c.Log.NewLogger()
	l.Info("config is read", slog.String("path", p))

	ctx.Context = Into(ctx.Context, c)
	ctx.Context = log.Into(ctx.Context, l)

	return c, nil
}
