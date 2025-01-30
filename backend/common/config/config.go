package config

import (
	"flag"
	"github.com/bsthun/gut"
	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Version     *string
	WebListen   []*string `env:"STRAWHOUSE_WEB_LISTEN" yaml:"webListen" validate:"required"`
	ProtoListen []*string `env:"STRAWHOUSE_PROTO_LISTEN" yaml:"protoListen" validate:"required"`
	DataRoot    *string   `env:"STRAWHOUSE_DATA_ROOT" yaml:"dataRoot" validate:"dirpath"`
	PogrebPath  *string   `env:"STRAWHOUSE_POGREB_PATH" yaml:"pogrebPath" validate:"dirpath"`
	PluginPath  *string   `env:"STRAWHOUSE_PLUGIN_PATH" yaml:"pluginPath" validate:"dirpath"`
	Key         *string   `env:"STRAWHOUSE_KEY" yaml:"key" validate:"required"`
	EnableXattr *bool     `env:"STRAWHOUSE_ENABLE_XATTR" yaml:"enableXattr" validate:"required"`
}

func Init() *Config {
	// * Parse arguments
	path := flag.String("config", "", "Path to config file")
	flag.Parse()

	// * Declare struct
	config := new(Config)

	// * Fallback to env
	if *path == "" {
		// * Parse env
		if err := env.Parse(config); err != nil {
			gut.Fatal("unable to parse environment variables", err)
		}
	} else {
		// * Read config
		yml, err := os.ReadFile(*path)
		if err != nil {
			gut.Fatal("unable to read configuration file", err)
		}

		// * Parse config
		if err := yaml.Unmarshal(yml, config); err != nil {
			gut.Fatal("unable to parse configuration file", err)
		}
	}

	// * Validate config
	if err := gut.Validate(config); err != nil {
		gut.Fatal("invalid configuration", err)
	}

	// * Normalize config
	*config.DataRoot, _ = filepath.Abs(*config.DataRoot)
	if err := os.MkdirAll(*config.DataRoot, 0700); err != nil {
		gut.Fatal("unable to create data root", err)
	}
	*config.PluginPath, _ = filepath.Abs(*config.PluginPath)
	if err := os.MkdirAll(*config.PluginPath, 0700); err != nil {
		gut.Fatal("unable to create plugin path", err)
	}

	return config
}
