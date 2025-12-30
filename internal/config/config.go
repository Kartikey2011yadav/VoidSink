package config

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	LogLevel  string `koanf:"log_level"`
	LogFile   string `koanf:"log_file"`
	LogFormat string `koanf:"log_format"`
	Metrics   struct {
		Enabled bool   `koanf:"enabled"`
		Addr    string `koanf:"addr"`
	} `koanf:"metrics"`
	Traps struct {
		HTTPInfinite struct {
			Enabled    bool   `koanf:"enabled"`
			Addr       string `koanf:"addr"`
			ServerName string `koanf:"server_name"`
		} `koanf:"http_infinite"`
		JSONInfinite struct {
			Enabled    bool   `koanf:"enabled"`
			Addr       string `koanf:"addr"`
			ServerName string `koanf:"server_name"`
		} `koanf:"json_infinite"`
	} `koanf:"traps"`
}

var k = koanf.New(".")

func Load(path string) (*Config, error) {
	// Determine parser based on extension
	var parser koanf.Parser
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		parser = yaml.Parser()
	case ".toml":
		parser = toml.Parser()
	default:
		parser = yaml.Parser()
	}

	// Load from file
	if err := k.Load(file.Provider(path), parser); err != nil {
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
