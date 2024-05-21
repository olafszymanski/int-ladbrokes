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
	App struct {
		LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
		Port     int    `env:"PORT" envDefault:"8080"`
	}
	Classes struct {
		RequestTimeout  time.Duration `env:"CLASSES_REQUEST_TIMEOUT" envDefault:"2s"`
		RequestInterval time.Duration `env:"CLASSES_REQUEST_INTERVAL" envDefault:"5s"`
	}
	Live struct {
		RequestTimeout  time.Duration `env:"LIVE_REQUEST_TIMEOUT" envDefault:"2s"`
		RequestInterval time.Duration `env:"LIVE_REQUEST_INTERVAL" envDefault:"2500ms"`
	}
	PreMatch struct {
		RequestTimeout  time.Duration `env:"PRE_MATCH_REQUEST_TIMEOUT" envDefault:"2s"`
		RequestInterval time.Duration `env:"PRE_MATCH_REQUEST_INTERVAL" envDefault:"10s"`
	}
	Storage struct {
		Address  string `env:"STORAGE_ADDRESS" envDefault:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
	}
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
