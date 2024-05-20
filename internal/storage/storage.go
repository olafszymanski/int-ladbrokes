package storage

import (
	"context"

	"github.com/olafszymanski/int-sdk/integration/pb"
	sdkStorage "github.com/olafszymanski/int-sdk/storage"
	"google.golang.org/protobuf/proto"
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
	return s.storage.Store(ctx, key, classes, 0)
}

func (s *Storage) GetEvent(ctx context.Context, hash, id string) (*pb.Event, error) {
	raw, err := s.storage.GetHashField(ctx, hash, id)
	if err != nil {
		return nil, err
	}
	ev := &pb.Event{}
	if err := proto.Unmarshal(raw, ev); err != nil {
		return nil, err
	}
	return ev, nil
}

func (s *Storage) StoreEvent(ctx context.Context, hash string, event *pb.Event) error {
	raw, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	return s.storage.StoreHashField(ctx, hash, event.Id, raw)
}

func (s *Storage) GetEvents(ctx context.Context, hash string) ([]*pb.Event, error) {
	raw, err := s.storage.GetHashFields(ctx, hash)
	if err != nil {
		return nil, err
	}
	ev := make([]*pb.Event, 0, len(raw))
	for _, r := range raw {
		e := &pb.Event{}
		if err := proto.Unmarshal(r, e); err != nil {
			return nil, err
		}
		ev = append(ev, e)
	}
	return ev, nil
}

func (s *Storage) StoreEvents(ctx context.Context, hash string, events []*pb.Event) error {
	rawEvs := make(map[string][]byte, len(events))
	for _, e := range events {
		raw, err := proto.Marshal(e)
		if err != nil {
			return err
		}
		rawEvs[e.Id] = raw
	}
	return s.storage.StoreHashFields(ctx, hash, rawEvs)
}

func (s *Storage) StoreNewEvents(ctx context.Context, hash string, events []*pb.Event) error {
	currIds, err := s.GetEventsIds(ctx, hash)
	if err != nil {
		return err
	}
	newEvents := getNewEvents(events, currIds)
	if len(newEvents) == 0 {
		return nil
	}
	return s.StoreEvents(ctx, hash, newEvents)
}

func (s *Storage) DeleteEvents(ctx context.Context, hash string, ids []string) error {
	return s.storage.DeleteHashFields(ctx, hash, ids)
}

func (s *Storage) GetEventsIds(ctx context.Context, hash string) ([]string, error) {
	return s.storage.GetHashFieldKeys(ctx, hash)
}

func (s *Storage) RemoveMissingEvents(ctx context.Context, hash string, events []*pb.Event) error {
	currIds, err := s.GetEventsIds(ctx, hash)
	if err != nil {
		return err
	}
	missIds := getMissingEventsIds(events, currIds)
	if len(missIds) == 0 {
		return nil
	}
	return s.DeleteEvents(ctx, hash, missIds)
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
