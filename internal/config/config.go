package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	LogLevel string `koanf:"log_level"`
	Traps    struct {
		HTTPInfinite struct {
			Enabled bool   `koanf:"enabled"`
			Addr    string `koanf:"addr"`
		} `koanf:"http_infinite"`
	} `koanf:"traps"`
}

var k = koanf.New(".")

func Load(path string) (*Config, error) {
	// Load from file
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		log.Printf("Error loading config file: %v", err)
	}

	// Load from environment variables
	k.Load(env.Provider("VOIDSINK_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "VOIDSINK_")), "_", ".", -1)
	}), nil)

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
