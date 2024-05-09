package poller

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/transform"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	eventsUrl = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/EventToOutcomeForClass/%s?simpleFilter=event.startTime:greaterThanOrEqual:%s&translationLang=en&responseFormat=json&prune=event&prune=market&childCount=event"

	lessThanFilter = "simpleFilter=event.startTime:lessThan:%s"

	liveEventsStorageKey     = "LIVE_EVENTS_%s"
	preMatchEventsStorageKey = "PRE_MATCH_EVENTS_%s"
)

var ErrMarshalEvent = fmt.Errorf("marshaling event failed")

type timeRange struct {
	start time.Duration
	end   time.Duration
}

var timeRanges = []timeRange{
	{-4 * time.Hour, 0},
	{0, 8 * time.Hour},
	{8 * time.Hour, 32 * time.Hour},
	{32 * time.Hour, 0}, // last element will not have end time
}

func (p *Poller) pollEvents(ctx context.Context, logger *logrus.Entry, sportType pb.SportType) error {
	var (
		startTime           time.Time
		liveEventsMapCh     = make(chan map[string][]byte)
		preMatchEventsMapCh = make(chan map[string][]byte)
	)
	defer func() {
		close(liveEventsMapCh)
		close(preMatchEventsMapCh)
	}()

	for {
		cls, err := p.storage.Get(ctx, fmt.Sprintf(classesStorageKey, sportType))
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			return err
		}
		if len(cls) == 0 {
			logger.Warn("skipping polling events, classes not found")
		} else {
			startTime = time.Now()
			go func() {
				evs, err := p.fetchEvents(cls)
				if err != nil {
					logger.WithError(err).Error("polling events failed")
					return
				}
				if len(evs) == 0 {
					logger.Warn("no events polled")
					return
				}
				liEvsMap, pmEvsMap, err := getEventsMaps(evs)
				if err != nil {
					logger.WithError(err).Error("converting events to maps failed")
					return
				}
				liveEventsMapCh <- liEvsMap
				preMatchEventsMapCh <- pmEvsMap
			}()
		}

		select {
		case evs := <-liveEventsMapCh:
			logger.WithField("events_length", len(evs)).Debug("live events polled")
			if err := p.storage.StoreHashFields(ctx, fmt.Sprintf(liveEventsStorageKey, sportType), evs); err != nil {
				return err
			}
			<-time.After(p.config.Events.RequestInterval - time.Since(startTime))
		case evs := <-preMatchEventsMapCh:
			logger.WithField("events_length", len(evs)).Debug("pre-match events polled")
			if err := p.storage.StoreHashFields(ctx, fmt.Sprintf(preMatchEventsStorageKey, sportType), evs); err != nil {
				return err
			}
			<-time.After(p.config.Events.RequestInterval - time.Since(startTime))
		case <-time.After(p.config.Events.RequestInterval):
			logger.Warn("events polling took longer than expected")
		}
	}
}

func (p *Poller) fetchEvents(classes []byte) ([]*pb.Event, error) {
	var (
		wg            = sync.WaitGroup{}
		lock          = sync.Mutex{}
		events        = make([]*pb.Event, 0)
		errCh         = make(chan error)
		done          = make(chan struct{})
		requestsCount = len(timeRanges)
	)
	defer func() {
		close(errCh)
		close(done)
	}()

	wg.Add(requestsCount)
	for i := 0; i < requestsCount; i++ {
		i := i
		go func() {
			defer wg.Done()

			u := getEventsUrl(
				classes,
				timeRanges,
				i,
				requestsCount,
			)
			evs, err := p.getEvents(u, p.config.Events.RequestTimeout)
			if err != nil {
				errCh <- err
				return
			}
			lock.Lock()
			events = append(events, evs...)
			lock.Unlock()
		}()
	}
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return events, nil
	case err := <-errCh:
		return nil, err
	}
}

func (p *Poller) getEvents(url string, timeout time.Duration) ([]*pb.Event, error) {
	res, err := p.httpClient.Get(url, httptls.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transform.TransformEvents(res.Body)
}

func getRequestTimes(timeRanges []timeRange, iteration int) (time.Time, time.Time) {
	n := time.Now().UTC()
	return n.Add(timeRanges[iteration].start), n.Add(timeRanges[iteration].end)
}

func getEventsUrl(classes []byte, timeRanges []timeRange, iteration, requestsCount int) string {
	st, et := getRequestTimes(timeRanges, iteration)
	if iteration == requestsCount-1 {
		return fmt.Sprintf(eventsUrl, classes, st.Format(time.RFC3339))
	}
	return fmt.Sprintf(
		fmt.Sprintf("%s&%s", eventsUrl, lessThanFilter),
		classes,
		st.Format(time.RFC3339),
		et.Format(time.RFC3339),
	)
}

// returns two maps: one for live events and one for pre-match events
func getEventsMaps(events []*pb.Event) (map[string][]byte, map[string][]byte, error) {
	li, pm := make(map[string][]byte), make(map[string][]byte)
	for _, e := range events {
		b, err := proto.Marshal(e)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %s", ErrMarshalEvent, err)
		}
		if e.IsLive {
			li[*e.ExternalId] = b
		} else {
			pm[*e.ExternalId] = b
		}
	}
	return li, pm, nil
}
