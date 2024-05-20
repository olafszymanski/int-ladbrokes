package main

import (
	"context"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-ladbrokes/internal/poller"
	"github.com/olafszymanski/int-ladbrokes/internal/storage"
	"github.com/olafszymanski/int-sdk/http"
	"github.com/olafszymanski/int-sdk/integration/pb"
	sdkStorage "github.com/olafszymanski/int-sdk/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("starting service...")
	logrus.SetLevel(logrus.DebugLevel)

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
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

	select {}
}
