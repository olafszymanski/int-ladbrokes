package poller

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
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
					lock.Lock()
					pollInfo[id].polling = false
					lock.Unlock()
					return
				}

				ev, err := p.storage.GetEvent(ctx, hash, id)
				if err != nil {
					logger.WithField("event_external_id", id).WithError(err).Warn("event not found in storage, probably finished and removed")
					lock.Lock()
					pollInfo[id].polling = false
					lock.Unlock()
					return
				}
				if err := p.updateEvent(ctx, logger, update, hash, ev); err != nil {
					errCh <- fmt.Errorf("failed to update event: %s", err)
					return
				}
				// Shouldn't store when finished
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

func (p *Poller) updateEvent(ctx context.Context, logger *logrus.Entry, update *transform.Update, hash string, event *pb.Event) error {
	for t, update := range update.Data {
		for _, data := range update {
			switch t {
			case mapping.EventUpdateType:
				u, err := transform.UnmarshalUpdate[model.EventUpdate](data.RawData)
				if err != nil {
					return fmt.Errorf("failed to unmarshal event update: %s", err)
				}
				if err := p.handleEventUpdate(ctx, u, data.ID, hash, event); err != nil {
					return fmt.Errorf("failed to handle event update: %s", err)
				}
			case mapping.MarketUpdateType:
				u, err := transform.UnmarshalUpdate[model.MarketUpdate](data.RawData)
				if err != nil {
					return fmt.Errorf("failed to unmarshal market update: %s", err)
				}
				handleMarketUpdate(logger, u, data.ID, event)
			case mapping.SelectionUpdateType:
				u, err := transform.UnmarshalUpdate[model.SelectionUpdate](data.RawData)
				if err != nil {
					return fmt.Errorf("failed to unmarshal selection update: %s", err)
				}
				if err := handleSelectionUpdate(u, data.ID, event); err != nil {
					return fmt.Errorf("failed to handle selection update: %s", err)
				}
			case mapping.PriceUpdateType:
				u, err := transform.UnmarshalUpdate[model.PriceUpdate](data.RawData)
				if err != nil {
					return fmt.Errorf("failed to unmarshal price update: %s", err)
				}
				if err := handlePriceUpdate(u, data.ID, event); err != nil {
					return fmt.Errorf("failed to handle price update: %s", err)
				}
			}
		}
	}
	return nil
}

func (p *Poller) handleEventUpdate(ctx context.Context, update *model.EventUpdate, updateId, hash string, event *pb.Event) error {
	if !transform.IsEventFinished(update) {
		return nil
	}
	if err := p.storage.DeleteEvents(ctx, hash, []string{updateId}); err != nil { // TODO: Create DeleteEvent method
		return fmt.Errorf("failed to remove finished event (%s) from storage: %s", event.ExternalId, err)
	}
	return nil
}

func handleMarketUpdate(logger *logrus.Entry, update *model.MarketUpdate, updateId string, event *pb.Event) {
	found := false

	for i, m := range event.Markets {
		if updateId != m.ExternalId {
			continue
		}

		var (
			ok         bool
			marketType = mapping.MoneyLineMarketType
		)
		if update.GroupNames != nil {
			marketType, ok = update.GroupNames["en"]
			if !ok {
				logger.WithField("group_names", update.GroupNames).Error("market update group names map changed")
				continue
			}
		}

		if _, ok := mapping.MarketTypes[marketType]; !ok {
			// TODO: Add check for unknown market types
			logger.WithField("market_type", marketType).Warn("unknown market type")
			continue
		}

		if transform.IsMarketRemoved(update) {
			event.Markets = append(event.Markets[:i], event.Markets[i+1:]...)
			return
		}

		a := !transform.IsMarketSuspended(update)
		for _, o := range m.Outcomes {
			o.IsAvailable = a
		}

		found = true
	}

	if !found {
		var (
			name           = ""
			marketType, ok = mapping.MarketTypes[update.GroupNames["en"]]
		)
		if !ok {
			// TODO: Add check for unknown market types
			logger.WithField("group_names", update.GroupNames).Warn("unknown market type")
			return
		}
		if marketType == pb.MarketType_PLAYER_TOTAL_POINTS || marketType == pb.MarketType_PLAYER_TOTAL_ASSISTS ||
			marketType == pb.MarketType_PLAYER_TOTAL_REBOUNDS || marketType == pb.MarketType_PLAYER_TOTAL_3_POINTERS {
			name = strings.Split(update.Names["en"], " (")[0]
		}

		mk := &pb.Market{
			Type:       marketType,
			ExternalId: updateId,
			Outcomes:   []*pb.Outcome{},
		}
		if name != "" {
			mk.Name = &name
		}
		event.Markets = append(event.Markets, mk)
	}
}

func handleSelectionUpdate(update *model.SelectionUpdate, updateId string, event *pb.Event) error {
	// b, _ := json.MarshalIndent(update, "", "  ")
	// fmt.Println("selec update for", update.EvMktID, updateId, string(b))
	for _, m := range event.Markets {
		if fmt.Sprint(update.EvMktID) != m.ExternalId {
			continue
		}

		if len(m.Outcomes) == 0 {
			n := update.Names["en"]
			m.Outcomes = append(m.Outcomes, &pb.Outcome{
				Type:        pb.Outcome_COMPETITOR,
				ExternalId:  updateId,
				Odds:        &pb.Odds{},
				Name:        &n,
				IsAvailable: !transform.IsSelectionSuspended(update),
			})
		}
		for _, outcome := range m.Outcomes {
			if updateId != outcome.ExternalId {
				continue
			}
			outcome.IsAvailable = !transform.IsSelectionSuspended(update)
			return updateOdds(outcome, update.LpNum, update.LpDen)
		}
	}
	return nil
}

func handlePriceUpdate(update *model.PriceUpdate, updateId string, event *pb.Event) error {
	for _, m := range event.Markets {
		for _, outcome := range m.Outcomes {
			if updateId != outcome.ExternalId {
				continue
			}
			return updateOdds(outcome, update.LpNum, update.LpDen)
		}
	}
	return nil
}

func updateOdds(outcome *pb.Outcome, numerator, denominator string) error {
	outcome.Odds.Numerator = numerator
	outcome.Odds.Denominator = denominator

	num, den, err := getFractionalOdds(numerator, denominator)
	if err != nil {
		return fmt.Errorf("failed to get fractional odds from update: %s", err)
	}
	res := num / den

	outcome.Odds.Decimal = math.Floor(res*100)/100 + 1
	if res > 1 {
		outcome.Odds.American = fmt.Sprint(math.Round(res * 100))
		return nil
	}
	outcome.Odds.American = fmt.Sprint(math.Round(-100 / res))
	return nil
}

func getFractionalOdds(numerator, denominator string) (float64, float64, error) {
	num, err := strconv.ParseFloat(numerator, 64)
	if err != nil {
		return 0, 0, err
	}
	den, err := strconv.ParseFloat(denominator, 64)
	if err != nil {
		return 0, 0, err
	}
	return num, den, nil
}

func getRequestBody(id string, requestBodyParts []string) []byte {
	return []byte(fmt.Sprintf("CL0000S0001sEVENT0%[1]s!!!!%[2]sS0001SEVENT0%[1]s!!!!%[3]s", id, requestBodyParts[0], requestBodyParts[1]))
}
