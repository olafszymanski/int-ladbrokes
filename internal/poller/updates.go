package poller

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-ladbrokes/internal/transform"
	sdkHttp "github.com/olafszymanski/int-sdk/http"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	updatesUrl = "https://push-lcm.ladbrokes.com/push"

	defaultRequestBody = "CL0000S0002sEVENT0%[1]sSEVENT0%[1]s!!!!!!!!!0"
)

type pollingInfo struct {
	polling bool
	body    []byte
}

func (p *Poller) pollUpdates(ctx context.Context, logger *logrus.Entry, sportType pb.SportType) error {
	logrus.Debug("polling updates")

	var (
		lock     sync.RWMutex
		pollInfo = make(map[string]*pollingInfo)
		errCh    = make(chan error)
	)
	defer close(errCh)

	go func() {
		for {
			ids, err := p.storage.GetHashFieldKeys(ctx, fmt.Sprintf(liveEventsStorageKey, sportType))
			if err != nil {
				errCh <- err
				return
			}
			if len(ids) < 1 {
				continue
			}
			for _, id := range ids {
				id := id

				lock.Lock()
				if d, ok := pollInfo[id]; ok && d.polling {
					lock.Unlock()
					continue
				} else if !ok {
					pollInfo[id] = &pollingInfo{
						polling: false,
						body:    nil,
					}
				}

				go func() {
					pollInfo[id].polling = true
					body := pollInfo[id].body
					if pollInfo[id].body == nil {
						body = []byte(fmt.Sprintf(defaultRequestBody, id))
					}
					lock.Unlock()

					ti := time.Now()
					update, err := p.getUpdates(body, time.Second*60)
					if err != nil {
						errCh <- fmt.Errorf("failed to receive update: %s", err)
						return
					}
					ev, err := p.getEventFromStorage(ctx, sportType, id)
					if err != nil {
						errCh <- fmt.Errorf("failed to get event from storage: %s", err)
						return
					}
					if err := updateEvent(update, ev); err != nil {
						errCh <- fmt.Errorf("failed to update event: %s", err)
						return
					}
					if err := p.storeEventInStorage(ctx, sportType, id, ev); err != nil {
						errCh <- fmt.Errorf("failed to store event: %s", err)
						return
					}

					lock.Lock()
					pollInfo[id].polling = false
					pollInfo[id].body = getRequestBody(id, update.RequestBodyParts)
					lock.Unlock()

					logger.WithFields(logrus.Fields{
						"start_time": ti,
						"end_time":   time.Now(),
					}).Debug("raw event received")
				}()
			}
		}
	}()
	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

func (p *Poller) getUpdates(requestBody []byte, timeout time.Duration) (*transform.Update, error) {
	res, err := p.httpClient.Do(&sdkHttp.Request{
		Method:  http.MethodPost,
		URL:     updatesUrl,
		Body:    requestBody,
		Timeout: timeout,
	})
	if err != nil {
		return nil, err
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transform.TransformUpdates(res.Body)
}

func (p *Poller) getEventFromStorage(ctx context.Context, sportType pb.SportType, id string) (*pb.Event, error) {
	raw, err := p.storage.GetHashField(ctx, fmt.Sprintf(liveEventsStorageKey, sportType), id)
	if err != nil {
		return nil, err
	}
	ev := &pb.Event{}
	if err := proto.Unmarshal(raw, ev); err != nil {
		return nil, err
	}
	return ev, nil
}

func (p *Poller) storeEventInStorage(ctx context.Context, sportType pb.SportType, id string, ev *pb.Event) error {
	raw, err := proto.Marshal(ev)
	if err != nil {
		return err
	}
	return p.storage.StoreHashField(ctx, fmt.Sprintf(liveEventsStorageKey, sportType), id, raw)
}

func updateEvent(update *transform.Update, event *pb.Event) error {
	for t, update := range update.Data {
		for _, data := range update {
			switch t {
			// case mapping.EventUpdateType:
			// 	u, err := transform.UnmarshalUpdate[model.EventUpdate](data.RawData)
			// 	if err != nil {
			// 		return fmt.Errorf("failed to unmarshal update: %s", err)
			// 	}
			case mapping.PriceUpdateType:
				u, err := transform.UnmarshalUpdate[model.PriceUpdate](data.RawData)
				if err != nil {
					return fmt.Errorf("failed to unmarshal update: %s", err)
				}
				for _, m := range event.Markets {
					for _, o := range m.Outcomes {
						if data.ID == o.ExternalId {
							o.Odds.Numerator = u.LpNum
							o.Odds.Denominator = u.LpDen
							// TODO: Add more
						}
					}
				}
			}
		}
	}
	return nil
}

func getRequestBody(id string, requestBodyParts []string) []byte {
	return []byte(fmt.Sprintf("CL0000S0001sEVENT0%[1]s!!!!%[2]sS0001SEVENT0%[1]s!!!!%[3]s", id, requestBodyParts[0], requestBodyParts[1]))
}
