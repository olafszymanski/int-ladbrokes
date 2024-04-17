package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Classes struct {
		RequestTimeout  time.Duration `env:"CLASSES_REQUEST_TIMEOUT" envDefault:"2s"`
		RequestInterval time.Duration `env:"CLASSES_REQUEST_INTERVAL" envDefault:"5s"`
	}
	PreMatch struct {
		RequestTimeout  time.Duration `env:"PRE_MATCH_REQUEST_TIMEOUT" envDefault:"10s"`
		RequestInterval time.Duration `env:"PRE_MATCH_REQUEST_INTERVAL" envDefault:"10s"`
	}
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
