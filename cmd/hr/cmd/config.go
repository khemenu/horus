package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/yaml.v2"
	"khepri.dev/horus"
	"khepri.dev/horus/cmd/conf"
	"khepri.dev/horus/ent"
	"khepri.dev/horus/internal/fx"
)

type Config struct {
	path string

	Grpc GrpcConfig `yaml:"grpc"`

	Db       DbConfig       `yaml:"db"`
	Client   ClientConfig   `yaml:"client"`
	Reporter ReporterConfig `yaml:"reporter"`

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

	has_token bool
	has_actor bool

	db     *ent.Client
	server *buff_server
	client horus.Client
}

func (c *ClientConfig) isBareServer() bool {
	return !(c.has_token || c.has_actor)
}

func (c *ClientConfig) toBeBareServe() error {
	if !c.isBareServer() {
		return fmt.Errorf("this operation requires server to be bare server")
	}

	return nil
}

func (c *ClientConfig) notToBeBareServe() error {
	if c.isBareServer() {
		return fmt.Errorf("this operation cannot be run on the bare server; please provide a token or an actor")
	}

	return nil
}

type ReporterConfig struct {
	Format   string `yaml:"format"` // <"plain"> | "template" | "json"
	Template string `yaml:"template"`

	reporter Reporter
}

func (c *ReporterConfig) Report(v any, plain string) (err error) {
	o := plain
	if c.reporter != nil {
		o, err = c.reporter.Report(v)
		if err != nil {
			return err
		}
	}

	fmt.Println(o)
	return nil
}

func (c *ReporterConfig) ExitWithErr(err error) cli.ExitCoder {
	if s, ok := status.FromError(err); ok {
		return cli.Exit(err.Error(), int(s.Code()))
	}
	if ent.IsNotFound(err) {
		return cli.Exit(err.Error(), int(codes.NotFound))
	}

	return cli.Exit(err.Error(), 1)
}

type ClientTargetConfig struct {
	Schema  string `yaml:"schema"`
	Address string `yaml:"address"`
}

type ClientDbConfig struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`

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
	fx.Default(&c.Reporter.Format, "plain")
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

var CmdGetConfig = &cli.Command{
	Name: "config",
	Action: func(ctx *cli.Context) error {
		conf := ConfFrom(ctx.Context)

		o, err := yaml.Marshal(conf)
		if err != nil {
			panic(fmt.Errorf("marshal config into YAML: %w", err))
		}

		return conf.Reporter.Report(conf, string(o))
	},
}
