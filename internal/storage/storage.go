package storage

import (
	"context"
	"encoding/json"

	"github.com/olafszymanski/int-sdk/integration/pb"
	sdkStorage "github.com/olafszymanski/int-sdk/storage"
)

type Storage struct {
	storage sdkStorage.Storager
}

func NewStorage(storage sdkStorage.Storager) *Storage {
	return &Storage{
		storage: storage,
	}
}

func (s *Storage) GetClasses(ctx context.Context, key string) ([]byte, error) {
	cls, err := s.storage.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return cls, nil
}

func (s *Storage) StoreClasses(ctx context.Context, key string, classes []byte) error {
	return s.storage.Set(ctx, key, classes, 0)
}

func (s *Storage) GetEvent(ctx context.Context, hash, id string) (*pb.Event, error) {
	raw, err := s.storage.GetMapValue(ctx, hash, id)
	if err != nil {
		return nil, err
	}

	var ev pb.Event
	if err := json.Unmarshal(raw, &ev); err != nil {
		return nil, err
	}
	return &ev, nil
}

func (s *Storage) StoreEvent(ctx context.Context, hash string, event *pb.Event) error {
	raw, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return s.storage.SetMapValue(ctx, hash, event.ExternalId, raw)
}

func (s *Storage) GetEvents(ctx context.Context, hash string) ([]*pb.Event, error) {
	raw, err := s.storage.GetMapValues(ctx, hash)
	if err != nil {
		return nil, err
	}

	evs := make([]*pb.Event, 0, len(raw))
	for _, r := range raw {
		var ev pb.Event
		if err := json.Unmarshal(r, &ev); err != nil {
			return nil, err
		}
		evs = append(evs, &ev)
	}
	return evs, nil
}

func (s *Storage) StoreEvents(ctx context.Context, hash string, events []*pb.Event) error {
	rawEvs := make(map[string]any, len(events))
	for _, e := range events {
		raw, err := json.Marshal(e)
		if err != nil {
			return err
		}
		rawEvs[e.ExternalId] = raw
	}
	return s.storage.SetMapValues(ctx, hash, rawEvs)
}

func (s *Storage) StoreNewEvents(ctx context.Context, hash string, events []*pb.Event) error {
	ids, err := s.GetEventsIds(ctx, hash)
	if err != nil {
		return err
	}
	new := getNewEvents(events, ids)
	if len(new) == 0 {
		return nil
	}
	return s.StoreEvents(ctx, hash, new)
}

func (s *Storage) DeleteEvents(ctx context.Context, hash string, ids []string) error {
	if err := s.storage.DeleteMapKeys(ctx, hash, ids); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEventsIds(ctx context.Context, hash string) ([]string, error) {
	return s.storage.GetMapKeys(ctx, hash)
}

func (s *Storage) RemoveMissingEvents(ctx context.Context, hash string, events []*pb.Event) error {
	curr, err := s.GetEventsIds(ctx, hash)
	if err != nil {
		return err
	}
	miss := getMissingEventsIds(events, curr)
	if len(miss) == 0 {
		return nil
	}
	return s.DeleteEvents(ctx, hash, miss)
}

func getNewEvents(events []*pb.Event, currentEventsIds []string) []*pb.Event {
	evs := make([]*pb.Event, 0)
	for _, event := range events {
		found := false
		for _, id := range currentEventsIds {
			if event.ExternalId == id {
				found = true
				break
			}
		}
		if !found {
			evs = append(evs, event)
		}
	}
	return evs
}

func getMissingEventsIds(events []*pb.Event, currentEventsIds []string) []string {
	evsIds := make(map[string]struct{}, len(events))
	for _, e := range events {
		evsIds[e.ExternalId] = struct{}{}
	}

	ids := make([]string, 0)
	for _, id := range currentEventsIds {
		if _, ok := evsIds[id]; !ok {
			ids = append(ids, id)
		}
	}
	return ids
}
