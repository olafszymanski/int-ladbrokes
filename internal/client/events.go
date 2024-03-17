package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/transformer"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

const (
	liveEventsUrl     = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/EventToOutcomeForClass/%s?simpleFilter=event.isStarted:isTrue&simpleFilter=event.startTime:greaterThanOrEqual:%s&translationLang=en&responseFormat=json&prune=event&prune=market&childCount=event"
	preMatchEventsUrl = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/EventToOutcomeForClass/%s?simpleFilter=event.isStarted:isFalse&simpleFilter=event.startTime:greaterThanOrEqual:%s&translationLang=en&responseFormat=json&prune=event&prune=market&childCount=event"
)

const lessThanFilter = "simpleFilter=event.startTime:lessThan:%s"

func (c *client) getEvents(url, classes string, maxConcurrentRequests, timeout int) ([]*pb.Event, error) {
	var (
		t     = time.Now().UTC()
		wg    = sync.WaitGroup{}
		evsCh = make(chan *pb.Event)
	)

	// Interval in which events are fetched, it's fibonacci sequence starting with 2
	in := 2
	for i := 0; i < maxConcurrentRequests; i++ {
		i := i

		wg.Add(1)
		go func() {
			defer wg.Done()

			in += in - 1

			u := getEventsUrl(url, classes, t, time.Duration(in)*time.Hour, i, maxConcurrentRequests)
			evs, err := c.fetchEvents(u, timeout)
			if err != nil {
				logrus.WithError(err).Error("fetching events failed")
				return
			}
			for _, e := range evs {
				evsCh <- e
			}
		}()
	}
	go func() {
		wg.Wait()
		close(evsCh)
	}()

	evs := make([]*pb.Event, 0, len(evsCh))
	for e := range evsCh {
		evs = append(evs, e)
	}
	return evs, nil
}

func (c *client) fetchEvents(url string, timeout int) ([]*pb.Event, error) {
	res, err := c.httpClient.Get(url, httptls.WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transformer.TransformEvents(res.Body)
}

func getEventsUrl(url, classes string, currentTime time.Time, interval time.Duration, currentIteration, maxConcurrentRequests int) string {
	currentTime = currentTime.Add(interval * time.Duration(currentIteration))

	fu := fmt.Sprintf("%s&%s", url, lessThanFilter)
	u := fmt.Sprintf(fu, classes, currentTime.Format(time.RFC3339), currentTime.Add(interval).Format(time.RFC3339))
	if currentIteration == maxConcurrentRequests-1 {
		u = fmt.Sprintf(url, classes, currentTime.Format(time.RFC3339))
	}
	return u
}
