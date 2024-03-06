package main

import (
	"context"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/client"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting service...")

	c := client.NewClient()

	ctx := context.Background()

	t := time.Now()
	for _, tp := range []pb.SportType{
		pb.SportType_BASKETBALL,
	} {
		_, err := c.GetPreMatch(ctx, &pb.Request{
			SportType: tp,
		})
		if err != nil {
			panic(err)
		}
		logrus.WithField("duration", time.Since(t)).Info("Pre match events fetched")
	}
}
