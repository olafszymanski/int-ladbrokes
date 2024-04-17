package poller

import (
	"fmt"
	"sync"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/transform"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

const (
	preMatchEventsUrl = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/EventToOutcomeForClass/%s?simpleFilter=event.isStarted:isFalse&simpleFilter=event.startTime:greaterThanOrEqual:%s&translationLang=en&responseFormat=json&prune=event&prune=market&childCount=event"
)

const lessThanFilter = "simpleFilter=event.startTime:lessThan:%s"

const concurrentRequests int = 4

var ErrMarshalEvent = fmt.Errorf("marshaling event failed")

var intervalPeriods = []time.Duration{
	1 * time.Hour,
	2 * time.Hour,
	3 * time.Hour,
	6 * time.Hour,
}

func (p *Poller) pollEvents(classes []byte) ([]*pb.Event, error) {
	var (
		ct     = time.Now().UTC()
		wg     = sync.WaitGroup{}
		lock   = sync.Mutex{}
		events = make([]*pb.Event, 0)
		errCh  = make(chan error)
		done   = make(chan struct{})
	)
	defer func() {
		close(errCh)
		close(done)
	}()

	wg.Add(concurrentRequests)
	for i := 0; i < concurrentRequests; i++ {
		i := i
		go func() {
			defer wg.Done()

			u := getEventsUrl(
				classes,
				ct,
				intervalPeriods[i],
				i,
			)
			evs, err := p.getEvents(u, p.config.PreMatch.RequestTimeout)
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
	res, err := p.httpClient.Get(
		url,
		httptls.WithTimeout(timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transform.TransformEvents(res.Body)
}

func getEventsUrl(classes []byte, currentTime time.Time, interval time.Duration, currentIteration int) string {
	currentTime = currentTime.Add(interval * time.Duration(currentIteration))

	fu := fmt.Sprintf("%s&%s", preMatchEventsUrl, lessThanFilter)
	u := fmt.Sprintf(fu, classes, currentTime.Format(time.RFC3339), currentTime.Add(interval).Format(time.RFC3339))
	if currentIteration == concurrentRequests-1 {
		u = fmt.Sprintf(preMatchEventsUrl, classes, currentTime.Format(time.RFC3339))
	}
	return u
}
