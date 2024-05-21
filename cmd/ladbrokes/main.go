package main

import (
	"context"

	"github.com/olafszymanski/int-ladbrokes/internal/client"
	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-ladbrokes/internal/poller"
	"github.com/olafszymanski/int-ladbrokes/internal/storage"
	"github.com/olafszymanski/int-sdk/http"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/integration/server"
	sdkStorage "github.com/olafszymanski/int-sdk/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("starting service...")

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	if cfg.App.LogLevel != "" {
		level, err := logrus.ParseLevel(cfg.App.LogLevel)
		if err != nil {
			logrus.WithError(err).Fatal("failed to parse log level")
		}
		logrus.SetLevel(level)
	}

	s := storage.NewStorage(sdkStorage.NewMemoryStorage())

	httpCl := http.NewClient()

	ctx := context.Background()

	go func() {
		p, err := poller.NewPoller(cfg, httpCl, s)
		if err != nil {
			logrus.WithError(err).Fatal("failed to create poller")
		}
		if err = p.Run(ctx, pb.SportType_BASKETBALL); err != nil {
			logrus.WithError(err).Fatal("failed to run poller")
		}
	}()

	cl := client.NewClient(cfg, httpCl, s)
	server.Start(cl, cfg.App.Port)
}
