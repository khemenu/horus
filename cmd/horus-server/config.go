package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
	"khepri.dev/horus"
	"khepri.dev/horus/provider"
	"khepri.dev/horus/server"
)

func defaultV[T comparable](target *T, v T) {
	var zero T
	if *target == zero {
		*target = v
	}
}

type Config struct {
	path string

	Rest ServerConfig `yaml:"rest"`
	Grpc ServerConfig `yaml:"grpc"`

	App AppConfig `yaml:"app"`
	Db  DbConfig  `yaml:"db"`

	Providers ProviderConfig `yaml:"providers"`

	Log   LogConfig   `yaml:"log"`
	Debug DebugConfig `yaml:"debug"`
}

func (c *Config) toHorusConf() *horus.Config {
	return &horus.Config{
		AppDomain: c.App.Domain,
		AppPrefix: c.App.Prefix,
	}
}

func (c *Config) toRestServerConf() *server.RestServerConfig {
	providers := []horus.Provider{}
	if c.Providers.GoogleOauth2 != nil {
		providers = append(providers, provider.GoogleOauth2(*c.Providers.GoogleOauth2))
	}

	return &server.RestServerConfig{
		Providers: providers,
		Debug: server.RestServerDebugConfig{
			Enabled:   c.Debug.Enabled,
			Unsecured: c.Debug.Unsecured,
		},
	}
}

func (c *Config) toGrpcServerConf() *server.GrpcServerConfig {
	return &server.GrpcServerConfig{}
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type AppConfig struct {
	Domain string `yaml:"domain"`
	Prefix string `yaml:"prefix"`
}

type DbConfig struct {
	Driver string `yaml:"driver"`
	Source string `yaml:"source"`
}

type ProviderConfig struct {
	GoogleOauth2 *provider.OauthProviderConfig `yaml:"google_oauth2"`
}

func (c *ProviderConfig) IsEmpty() bool {
	return c.GoogleOauth2 == nil
}

type LogConfig struct {
	Format string `yaml:"format"` // "text" | "json"
}

func (c *LogConfig) NewLogger() *slog.Logger {
	var logger *slog.Logger
	switch c.Format {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	default:
		panic("unreachable")
	}

	return logger
}

type DebugConfig struct {
	Enabled   bool `yaml:"enabled"`
	Unsecured bool `yaml:"unsecured"`
	UseMemDb  bool `yaml:"use_mem_db"`
}

func ParseArgs(args []string) (*Config, error) {
	conf := &Config{}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&conf.path, "conf", "horus.yaml", "path to a config file")
	flags.Parse(args[1:])

	data, err := os.ReadFile(conf.path)
	if err != nil {
		return nil, fmt.Errorf("read config at %s: %w", conf.path, err)
	}

	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("unmarshal config at %s: %w", conf.path, err)
	}

	defaultV(&conf.Rest.Host, "localhost")
	defaultV(&conf.Rest.Port, 20000)
	defaultV(&conf.Grpc.Host, "localhost")
	defaultV(&conf.Grpc.Port, 20001)
	defaultV(&conf.App.Prefix, "/auth")
	defaultV(&conf.Log.Format, "text")

	errs := []error{}
	if !slices.Contains([]string{"text", "json"}, conf.Log.Format) {
		errs = append(errs, fmt.Errorf(`log.format must be one of "text" or "json": %s`, conf.Log.Format))
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return conf, nil
}
