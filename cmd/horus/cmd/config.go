package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
	"khepri.dev/horus/cmd/conf"
	"khepri.dev/horus/internal/fx"
)

type Config struct {
	path string

	Grpc GrpcConfig `yaml:"grpc"`
	Http HttpConfig `yaml:"http"`

	Db DbConfig `yaml:"db"`

	Log   conf.LogConfig `yaml:"log"`
	Debug DebugConfig    `yaml:"debug"`
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

func ParseArgs(args []string) (*Config, error) {
	c := &Config{}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&c.path, "conf", "horus.yaml", "path to a config file")
	flags.Parse(args[1:])

	data, err := os.ReadFile(c.path)
	if err != nil {
		return nil, fmt.Errorf("read config at %s: %w", c.path, err)
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("unmarshal config at %s: %w", c.path, err)
	}

	fx.Default(&c.Grpc.Host, "localhost")
	fx.Default(&c.Grpc.Port, 35122)
	fx.Default(&c.Http.Host, "localhost")
	fx.Default(&c.Http.Port, 35123)
	fx.Default(&c.Log.Format, "text")

	errs := []error{}
	if !slices.Contains([]string{"text", "json"}, c.Log.Format) {
		errs = append(errs, fmt.Errorf(`log.format must be one of "text" or "json": %s`, c.Log.Format))
	}

	if len(errs) != 0 {
		return nil, errors.Join(errs...)
	}

	return c, nil
}
