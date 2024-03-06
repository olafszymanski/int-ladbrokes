package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
		res, err := c.GetPreMatch(ctx, &pb.Request{
			SportType: tp,
		})
		if err != nil {
			panic(err)
		}
		logrus.WithField("duration", time.Since(t)).Info("Pre match events fetched")
		if err = writeEventsToFile(res.Events, fmt.Sprintf("results/%s.json", tp)); err != nil {
			panic(err)
		}
	}
}

func writeEventsToFile(events []*pb.Event, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := json.MarshalIndent(events, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	return err
}
