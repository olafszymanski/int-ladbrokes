package poller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const preMatchEventsStorageKey = "PRE_MATCH_EVENTS_%s"

var (
	ErrRequest              = fmt.Errorf("request failed")
	ErrUnexpectedStatusCode = fmt.Errorf("unexpected status code")
)

type Poller struct {
	config     *config.Config
	httpClient *httptls.HTTPClient
	storage    storage.Storager
}

func NewPoller(config *config.Config, httpClient *httptls.HTTPClient, storage storage.Storager) *Poller {
	return &Poller{
		config:     config,
		httpClient: httpClient,
		storage:    storage,
	}
}

func (p *Poller) Run(ctx context.Context, sportType pb.SportType) error {
	var (
		errCh  = make(chan error)
		logger = logrus.WithField("sport_type", sportType)
	)
	defer close(errCh)

	go func() {
		var (
			classesCh = make(chan []byte)
			startTime time.Time
		)
		defer close(classesCh)
		for {
			startTime = time.Now()
			go func() {
				cls, err := p.pollClasses(sportType)
				if err != nil {
					logger.WithError(err).Error("fetching classes failed")
				}
				if len(cls) == 0 {
					logger.Warn("no pre-match events fetched")
					return
				}
				classesCh <- cls
			}()

			select {
			case cls := <-classesCh:
				logger.WithField("classes_length", len(cls)).Debug("classes fetched")
				if err := p.storage.Store(ctx, fmt.Sprintf(classesStorageKey, sportType), cls, 0); err != nil {
					errCh <- err
					return
				}
				<-time.After(p.config.Classes.RequestInterval - time.Since(startTime))
			case <-time.After(p.config.Classes.RequestInterval):
				logger.Warn("classes polling took longer than expected")
			}
		}
	}()

	go func() {
		var (
			eventsMapCh = make(chan map[string]any)
			startTime   time.Time
		)
		defer close(eventsMapCh)
		for {
			cls, err := p.storage.Get(ctx, fmt.Sprintf(classesStorageKey, sportType))
			if err != nil && !errors.Is(err, storage.ErrNotFound) {
				errCh <- err
				return
			}
			if len(cls) == 0 {
				logger.WithError(err).Error("skipping fetching pre-match events, classes not found")
			} else {
				startTime = time.Now()
				go func() {
					evs, err := p.pollEvents(cls)
					if err != nil {
						logger.WithError(err).Error("fetching pre-match events failed")
					}
					if len(evs) == 0 {
						logger.Warn("no pre-match events fetched")
						return
					}
					evsMap, err := getEventsMap(evs)
					if err != nil {
						errCh <- err
						return
					}
					eventsMapCh <- evsMap
				}()
			}

			select {
			case evs := <-eventsMapCh:
				logger.WithField("events_length", len(evs)).Debug("pre-match events fetched")
				if err := p.storage.StoreHash(ctx, fmt.Sprintf(preMatchEventsStorageKey, sportType), evs); err != nil {
					errCh <- err
					return
				}
				<-time.After(p.config.PreMatch.RequestInterval - time.Since(startTime))
			case <-time.After(p.config.PreMatch.RequestInterval):
				logger.Warn("pre-match events polling took longer than expected")
			}
		}
	}()

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

func getEventsMap(events []*pb.Event) (map[string]any, error) {
	evs := make(map[string]any)
	for _, e := range events {
		b, err := proto.Marshal(e)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrMarshalEvent, err)
		}
		// TODO: Later change it to use e.Id
		evs[*e.ExternalId] = b
	}
	return evs, nil
}
