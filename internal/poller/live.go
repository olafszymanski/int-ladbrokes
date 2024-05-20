package poller

import (
	"context"
	"fmt"
	"maps"
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
			p.saveEvents(ctx, fmt.Sprintf(config.LiveEventsStorageKey, sportType), evs, func(currentEventsIds []string, marshaledEvents map[string][]byte) error {
				// we store only new live events as we don't want to override data coming from push updates
				e := filterNewEvents(currentEventsIds, marshaledEvents)
				if len(e) == 0 {
					return nil
				}
				return p.storage.StoreHashFields(ctx, hash, e)
			})
			<-time.After(p.config.Live.RequestInterval - time.Since(startTime))
		case <-time.After(p.config.Live.RequestInterval):
			logger.Warn("live events polling took longer than expected")
		}
	}
}

func filterNewEvents(ids []string, events map[string][]byte) map[string][]byte {
	e := make(map[string][]byte)
	maps.Copy(e, events)
	for k := range events {
		for _, id := range ids {
			if k == id {
				delete(e, k)
			}
		}
	}
	return e
}
