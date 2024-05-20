package poller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/transform"
	sdkHttp "github.com/olafszymanski/int-sdk/http"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
	"google.golang.org/protobuf/proto"
)

const (
	eventsUrl = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/EventToOutcomeForClass/%s?simpleFilter=event.startTime:greaterThanOrEqual:%s&translationLang=en&responseFormat=json&prune=event&prune=market&childCount=event"

	startTimelessThanFilter = "simpleFilter=event.startTime:lessThan:%s"
)

type saveFunc func(currentEventsIds []string, marshaledEvents map[string][]byte) error

func (p *Poller) pollEvents(ctx context.Context, baseUrl string, sportType pb.SportType, timeout time.Duration, timePeriods []timePeriod) ([]*pb.Event, error) {
	cls, err := p.storage.Get(ctx, fmt.Sprintf(classesStorageKey, sportType))
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, err
	}
	if len(cls) == 0 {
		return nil, nil
	}
	return p.fetchEvents(baseUrl, cls, timeout, timePeriods)
}

func (p *Poller) fetchEvents(baseUrl string, classes []byte, timeout time.Duration, timePeriods []timePeriod) ([]*pb.Event, error) {
	var (
		requestsCount = len(timePeriods)
		wg            = sync.WaitGroup{}
		lock          = sync.Mutex{}
		events        = make([]*pb.Event, 0)
		done          = make(chan struct{})
		errCh         = make(chan error)
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

			var url string
			if i == requestsCount-1 {
				url = getLastUrl(
					baseUrl,
					classes,
					&timePeriods[i],
				)
			} else {
				url = getUrl(
					baseUrl,
					classes,
					&timePeriods[i],
				)
			}

			evs, err := p.getEvents(url, timeout)
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
	res, err := p.httpClient.Do(&sdkHttp.Request{
		Method:  http.MethodGet,
		URL:     url,
		Timeout: timeout,
	})
	if err != nil {
		return nil, err
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transform.TransformEvents(res.Body)
}

func (p *Poller) saveEvents(ctx context.Context, hash string, events []*pb.Event, save saveFunc) error {
	mevs, err := marshalEvents(events)
	if err != nil {
		return err
	}
	ids, err := p.storage.GetHashFieldKeys(ctx, hash)
	if err != nil {
		return err
	}
	if err := p.removeMissingEvents(ctx, hash, ids, mevs); err != nil {
		return err
	}
	return save(ids, mevs)
}

func (p *Poller) removeMissingEvents(ctx context.Context, hash string, ids []string, events map[string][]byte) error {
	i := getMissingEventsIds(ids, events)
	if len(i) < 1 {
		return nil
	}
	return p.storage.DeleteHashFields(ctx, hash, i)
}

func getUrl(url string, classes []byte, timePeriod *timePeriod) string {
	st, et := timePeriod.getTimes()
	return fmt.Sprintf(
		fmt.Sprintf("%s&%s", url, startTimelessThanFilter),
		classes,
		st.Format(time.RFC3339),
		et.Format(time.RFC3339),
	)
}

func getLastUrl(url string, classes []byte, timePeriod *timePeriod) string {
	st, _ := timePeriod.getTimes()
	return fmt.Sprintf(url, classes, st.Format(time.RFC3339))
}

func getMissingEventsIds(ids []string, events map[string][]byte) []string {
	r := make([]string, 0)
	for _, k := range ids {
		if _, ok := events[k]; !ok {
			r = append(r, k)
		}
	}
	return r
}

func marshalEvents(events []*pb.Event) (map[string][]byte, error) {
	evs := make(map[string][]byte)
	for _, e := range events {
		b, err := proto.Marshal(e)
		if err != nil {
			return nil, err
		}
		evs[e.ExternalId] = b
	}
	return evs, nil
}
