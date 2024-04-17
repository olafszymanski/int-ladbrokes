package main

import (
	"context"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-ladbrokes/internal/poller"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("starting service...")

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.WithError(err).Fatal("failed to load config")
	}

	s := storage.NewMemoryStorage()

	cl := httptls.NewHTTPClient()

	ctx := context.Background()

	go func() {
		p := poller.NewPoller(cfg, cl, s)
		if err = p.Run(ctx, pb.SportType_BASKETBALL); err != nil {
			logrus.WithError(err).Fatal("failed to run poller")
		}
	}()

	<-ctx.Done()
}
