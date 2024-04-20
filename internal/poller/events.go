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
	eventsUrl = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/EventToOutcomeForClass/%s?simpleFilter=event.startTime:greaterThanOrEqual:%s&translationLang=en&responseFormat=json&prune=event&prune=market&childCount=event&referenceEachWayTerms=true"
)

const lessThanFilter = "simpleFilter=event.startTime:lessThan:%s"

var ErrMarshalEvent = fmt.Errorf("marshaling event failed")

var intervalPeriods = []time.Duration{
	4 * time.Hour,
	4 * time.Hour,
	8 * time.Hour,
	12 * time.Hour,
}

func (p *Poller) pollEvents(classes []byte) ([]*pb.Event, error) {
	var (
		wg            = sync.WaitGroup{}
		lock          = sync.Mutex{}
		events        = make([]*pb.Event, 0)
		errCh         = make(chan error)
		done          = make(chan struct{})
		currentTime   = time.Now().UTC()
		requestsCount = len(intervalPeriods)
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

			if i == 0 {
				currentTime = currentTime.Add(-intervalPeriods[i] * time.Hour)
			} else {
				currentTime = time.Now().UTC()
			}
			u := getEventsUrl(
				classes,
				currentTime,
				intervalPeriods[i],
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
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transform.TransformEvents(res.Body)
}

func getEventsUrl(classes []byte, currentTime time.Time, interval time.Duration, currentIteration, requestsCount int) string {
	currentTime = currentTime.Add(interval * time.Duration(currentIteration))

	fu := fmt.Sprintf("%s&%s", eventsUrl, lessThanFilter)
	u := fmt.Sprintf(fu, classes, currentTime.Format(time.RFC3339), currentTime.Add(interval).Format(time.RFC3339))
	if currentIteration == requestsCount-1 {
		u = fmt.Sprintf(eventsUrl, classes, currentTime.Format(time.RFC3339))
	}
	return u
}
