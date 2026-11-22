package config

import (
	"fmt"
	"log/slog"
	"os"
)

type section interface {
	fields() []fieldSpec
}

type defaulter interface {
	defaults()
}

type fieldSpec struct {
	name  string
	parse func(string) error
}

// loadFields populates each spec's target from its environment variable.
// Unset variables keep their default.
func loadFields(specs []fieldSpec) error {
	for _, s := range specs {
		v := os.Getenv(s.name)
		if v == "" {
			continue
		}

		if err := s.parse(v); err != nil {
			return fmt.Errorf("%s: %w", s.name, err)
		}
	}

	return nil
}

type ServerConfig struct {
	ListenAddr string
}

// fields returns list of fieldSpec, holding env definitions
// and pointers into the ServerConfig for env loading
func (cnf *ServerConfig) fields() []fieldSpec {
	return []fieldSpec{
		{"AMNOVA_LISTEN_ADDR", parseString(&cnf.ListenAddr)},
	}
}

func (cnf *ServerConfig) defaults() {
	cnf.ListenAddr = ":8080"
}

type LogConfig struct {
	Level slog.Level
}

// fields returns list of fieldSpec, holding env definitions
// and pointers into the LogConfig for env loading
func (cnf *LogConfig) fields() []fieldSpec {
	return []fieldSpec{
		{"AMNOVA_LOG_LEVEL", parseLevel(&cnf.Level)},
	}
}

func (cnf *LogConfig) defaults() {
	cnf.Level = slog.LevelInfo
}

type StorageConfig struct {
	DataDir string
}

// fields returns list of fieldSpec, holding env definitions
// and pointers into the StorageConfig for env loading
func (cnf *StorageConfig) fields() []fieldSpec {
	return []fieldSpec{
		{"AMNOVA_DATA_DIR", parseString(&cnf.DataDir)},
	}
}

func (cnf *StorageConfig) defaults() {
	cnf.DataDir = "data"
}

type Config struct {
	Server  ServerConfig
	Log     LogConfig
	Storage StorageConfig
}

func (cfg *Config) sections() []section {
	return []section{
		&cfg.Server,
		&cfg.Log,
		&cfg.Storage,
	}
}

// Load reads configuration from environment variables
// and returns a populated Config or an error.
func Load() (*Config, error) {
	cfg := &Config{}

	for _, s := range cfg.sections() {
		if d, ok := s.(defaulter); ok {
			d.defaults()
		}

		if err := loadFields(s.fields()); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func parseString(dst *string) func(string) error {
	return func(v string) error { *dst = v; return nil }
}

func parseLevel(dst *slog.Level) func(string) error {
	return func(v string) error {
		var l slog.Level
		if err := l.UnmarshalText([]byte(v)); err != nil {
			return fmt.Errorf("want debug, info, warn or error, got %q", v)
		}
		*dst = l
		return nil
	}
}
