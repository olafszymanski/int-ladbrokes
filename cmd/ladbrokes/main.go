package main

import (
	"github.com/olafszymanski/int-ladbrokes/internal/client"
	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-sdk/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting service...")

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	s := storage.NewMemoryStorage()

	client.NewClient(cfg, s)
}
