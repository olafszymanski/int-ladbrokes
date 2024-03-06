package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/client"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		in  string
		out string
	)
	flag.StringVar(&in, "input", "", "Input file path to use for dry run")
	flag.StringVar(&out, "output", "", "Save the output to a file")
	flag.Parse()
	if in == "" {
		panic("input filepath is required")
	}

	f, err := os.ReadFile(in)
	if err != nil {
		panic(err)
	}

	logrus.Info("Starting service...")

	c := client.NewTestClient(&client.TestClientConfig{
		Input: bytes.NewReader(f),
	})

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

		if out != "" {
			f, err := os.OpenFile(out, os.O_CREATE|os.O_RDWR, 0666)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if err = saveEvents(res.Events, f); err != nil {
				panic(err)
			}
		}
	}
}

func saveEvents(events []*pb.Event, file *os.File) error {
	b, err := json.MarshalIndent(events, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	return err
}
