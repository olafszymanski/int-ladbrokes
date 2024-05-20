package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

const (
	LiveEventsStorageKey     = "LIVE_EVENTS_%s"
	PreMatchEventsStorageKey = "PRE_MATCH_EVENTS_%s"
)

type Config struct {
	Classes struct {
		RequestTimeout  time.Duration `env:"CLASSES_REQUEST_TIMEOUT" envDefault:"2s"`
		RequestInterval time.Duration `env:"CLASSES_REQUEST_INTERVAL" envDefault:"5s"`
	}
	Events struct {
		RequestTimeout  time.Duration `env:"EVENTS_REQUEST_TIMEOUT" envDefault:"2s"`
		RequestInterval time.Duration `env:"EVENTS_REQUEST_INTERVAL" envDefault:"5s"`
	}
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
