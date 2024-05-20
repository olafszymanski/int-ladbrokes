package poller

import (
	"context"
	"fmt"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

const liveFilter = "simpleFilter=event.isStarted:isTrue"

func (p *Poller) pollLiveEvents(ctx context.Context, logger *logrus.Entry, sportType pb.SportType) error {
	var (
		startTime   time.Time
		timePeriods = []timePeriod{
			{-4 * time.Hour, 0}, // last (and first in this case) element does not have an end time
		}
		eventsCh = make(chan []*pb.Event)
	)
	defer close(eventsCh)

	for {
		startTime = time.Now()

		go func() {
			u := fmt.Sprintf("%s&%s", eventsUrl, liveFilter)
			evs, err := p.pollEvents(ctx, u, sportType, p.config.Live.RequestTimeout, timePeriods)
			if err != nil {
				logger.WithError(err).Error("polling live events failed")
				return
			}
			if len(evs) == 0 {
				logger.Warn("no live events polled")
				return
			}
			eventsCh <- evs
		}()

		select {
		case evs := <-eventsCh:
			logger.WithField("length", len(evs)).Debug("live events polled")

			hash := fmt.Sprintf(config.LiveEventsStorageKey, sportType)
			if err := p.storage.RemoveMissingEvents(ctx, hash, evs); err != nil {
				return err
			}
			if err := p.storage.StoreNewEvents(ctx, hash, evs); err != nil {
				return err
			}
			<-time.After(p.config.Live.RequestInterval - time.Since(startTime))
		case <-time.After(p.config.Live.RequestInterval):
			logger.Warn("live events polling took longer than expected")
		}
	}
}
