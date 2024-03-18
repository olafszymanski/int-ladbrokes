package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/client"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		liIn   string
		prIn   string
		out    string
		cpuOut string
	)
	flag.StringVar(&liIn, "live-input", "", "Live input file path to read the data from")
	flag.StringVar(&prIn, "pre-match-input", "", "Pre match input file path to read the data from")
	flag.StringVar(&out, "output", "", "Output file path to save the result to")
	flag.StringVar(&cpuOut, "cpu-output", "", "Output file path to save the CPU profile to")
	flag.Parse()
	if liIn == "" && prIn == "" {
		panic("input file path is required for either live or pre-match")
	}

	if cpuOut != "" {
		f, err := os.OpenFile(cpuOut, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	cfg := &client.TestClientConfig{}
	if liIn != "" {
		f, err := os.ReadFile(liIn)
		if err != nil {
			panic(err)
		}
		cfg.LiveInput = f
	}
	if prIn != "" {
		f, err := os.ReadFile(prIn)
		if err != nil {
			panic(err)
		}
		cfg.PreMatchInput = f
	}

	var (
		c   = client.NewTestClient(cfg)
		ctx = context.Background()
		t   = time.Now()
		wg  sync.WaitGroup
	)
	for _, tp := range []pb.SportType{
		pb.SportType_BASKETBALL,
	} {
		tp := tp

		wg.Add(2)
		go func() {
			defer wg.Done()
			if liIn != "" {
				res, err := c.GetLive(ctx, &pb.Request{
					SportType: tp,
				})
				if err != nil {
					panic(err)
				}

				if out != "" {
					f, err := os.OpenFile(
						getFilePath(out, "LIVE"),
						os.O_CREATE|os.O_RDWR,
						0666,
					)
					if err != nil {
						panic(err)
					}
					defer f.Close()

					if err = saveEvents(res.Events, f); err != nil {
						panic(err)
					}
				}
			}
		}()

		go func() {
			defer wg.Done()
			if prIn != "" {
				res, err := c.GetPreMatch(ctx, &pb.Request{
					SportType: tp,
				})
				if err != nil {
					panic(err)
				}

				if out != "" {
					f, err := os.OpenFile(
						getFilePath(out, "PRE_MATCH"),
						os.O_CREATE|os.O_RDWR,
						0666,
					)
					if err != nil {
						panic(err)
					}
					defer f.Close()

					if err = saveEvents(res.Events, f); err != nil {
						panic(err)
					}
				}
			}
		}()
		wg.Wait()

		logrus.WithField("duration", time.Since(t)).Info("events fetched")
	}
}

func getFilePath(filePath, prefix string) string {
	p := strings.Split(filePath, "/")
	p[len(p)-1] = fmt.Sprintf("%s_%s", prefix, p[len(p)-1])
	return strings.Join(p, "/")
}

func saveEvents(events []*pb.Event, file *os.File) error {
	b, err := json.MarshalIndent(events, "", "   ")
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	return err
}
