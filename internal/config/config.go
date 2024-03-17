package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	LiveEventsRequestTimeout            int `env:"LIVE_EVENTS_REQUEST_TIMEOUT" envDefault:"2000"`
	MaxLiveEventsConcurrentRequests     int `env:"MAX_LIVE_EVENTS_CONCURRENT_REQUESTS" envDefault:"4"`
	PreMatchEventsRequestTimeout        int `env:"PRE_MATCH_EVENTS_REQUEST_TIMEOUT" envDefault:"2000"`
	MaxPreMatchEventsConcurrentRequests int `env:"MAX_PRE_MATCH_EVENTS_CONCURRENT_REQUESTS" envDefault:"4"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cfg, nil
}
