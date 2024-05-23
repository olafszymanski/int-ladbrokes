package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

const preMatchFilter = "simpleFilter=event.isStarted:isFalse"

func (p *Poller) pollPreMatchEvents(ctx context.Context, logger *logrus.Entry, sportType pb.SportType) error {
	var (
		startTime   time.Time
		timePeriods = []timePeriod{
			{-4 * time.Hour, 0},
			{0, 8 * time.Hour},
			{8 * time.Hour, 32 * time.Hour},
			{32 * time.Hour, 0}, // last element does not have an end time
		}
		eventsCh = make(chan []*pb.Event)
	)
	defer close(eventsCh)

	for {
		startTime = time.Now()

		go func() {
			u := fmt.Sprintf("%s&%s", eventsUrl, preMatchFilter)
			evs, err := p.pollEvents(ctx, u, sportType, p.config.PreMatch.RequestTimeout, timePeriods)
			if err != nil {
				logger.WithError(err).Error("polling pre-match events failed")
				return
			}
			if len(evs) == 0 {
				logger.Warn("no pre-match events polled")
				return
			}
			eventsCh <- evs
		}()

		select {
		case evs := <-eventsCh:
			logger.WithField("length", len(evs)).Debug("pre-match events polled")

			hash := fmt.Sprintf(config.PreMatchEventsStorageKey, sportType)
			if err := p.storage.RemoveMissingEvents(ctx, hash, evs); err != nil {
				return fmt.Errorf("failed to remove missing pre-match events: %s", err)
			}
			if err := p.storage.StoreEvents(ctx, hash, evs); err != nil {
				return fmt.Errorf("failed to store pre-match events: %s", err)
			}
			<-time.After(p.config.PreMatch.RequestInterval - time.Since(startTime))
		case <-time.After(p.config.PreMatch.RequestInterval):
			logger.Warn("pre-match events polling took longer than expected")
		}
	}
}
