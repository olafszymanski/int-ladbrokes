package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	LiveEventsRequestTimeout     int `env:"LIVE_EVENTS_REQUEST_TIMEOUT" envDefault:"2000"`
	PreMatchEventsRequestTimeout int `env:"PRE_MATCH_EVENTS_REQUEST_TIMEOUT" envDefault:"2000"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
