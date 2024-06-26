package poller

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-ladbrokes/internal/transform"
	sdkHttp "github.com/olafszymanski/int-sdk/http"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
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
	logger.Debug("polling updates")

	var (
		lock     = sync.Mutex{}
		hash     = fmt.Sprintf(config.LiveEventsStorageKey, sportType)
		pollInfo = make(map[string]*pollingInfo)
		idCh     = make(chan string)
		errCh    = make(chan error)
	)
	defer close(errCh)

	go func() {
		for {
			ids, err := p.storage.GetEventsIds(ctx, hash)
			if err != nil {
				errCh <- fmt.Errorf("failed to get events ids for updates polling: %s", err)
				return
			}
			if len(ids) == 0 {
				continue
			}
			for _, id := range ids {
				id := id

				go func() {
					lock.Lock()
					defer lock.Unlock()
					if info, ok := pollInfo[id]; !ok {
						pollInfo[id] = &pollingInfo{
							polling: true,
							body:    nil,
						}
						idCh <- id
					} else if !info.polling {
						pollInfo[id].polling = true
						idCh <- id
					}
				}()
			}
		}
	}()

	for {
		select {
		case id := <-idCh:
			go func() {
				lock.Lock()
				body := pollInfo[id].body
				if body == nil {
					body = []byte(fmt.Sprintf(defaultRequestBody, id))
				}
				lock.Unlock()

				st := time.Now()

				update, err := p.getUpdates(body, time.Second*60)
				if err != nil {
					errCh <- fmt.Errorf("failed to receive update: %s", err)
					return
				}
				// no update received, we don't have to do anything
				if update == nil {
					logger.WithField("event_external_id", id).Debug("no update")
					return
				}

				ev, err := p.storage.GetEvent(ctx, hash, id)
				if err != nil {
					errCh <- fmt.Errorf("failed to get event from storage: %s", err)
					return
				}
				if err := updateEvent(update, ev); err != nil {
					errCh <- fmt.Errorf("failed to update event: %s", err)
					return
				}
				if err := p.storage.StoreEvent(ctx, hash, ev); err != nil {
					errCh <- fmt.Errorf("failed to save event: %s", err)
					return
				}

				lock.Lock()
				pollInfo[id].polling = false
				pollInfo[id].body = getRequestBody(id, update.RequestBodyParts)
				lock.Unlock()

				logger.WithFields(logrus.Fields{
					"event_external_id": id,
					"start_time":        st,
					"duration":          time.Since(st),
				}).Debug("update received")
			}()
			continue
		case err := <-errCh:
			return err
		}
	}
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
