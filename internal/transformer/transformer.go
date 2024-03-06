package transformer

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrDecodeResponse = fmt.Errorf("decoding response failed")
	ErrParseTime      = fmt.Errorf("parsing time failed")
	ErrParseMarket    = fmt.Errorf("parsing market failed")
)

type Transformer struct {
	logger *logrus.Entry
}

func NewTransformer(logger *logrus.Entry) *Transformer {
	return &Transformer{
		logger: logger,
	}
}

func (t *Transformer) UnmarshallClasses(reader io.Reader) ([]string, error) {
	var root model.ClassesRoot
	if err := json.NewDecoder(reader).Decode(&root); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDecodeResponse, err)
	}
	return t.transformClasses(&root), nil
}

func (t *Transformer) UnmarshallEvents(reader io.Reader) ([]*pb.Event, error) {
	var root model.EventsRoot
	if err := json.NewDecoder(reader).Decode(&root); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDecodeResponse, err)
	}
	return t.transformEvents(&root)
}

func (t *Transformer) transformClasses(classesRoot *model.ClassesRoot) []string {
	cls := make([]string, 0, len(classesRoot.SSResponse.Children))
	for _, c := range classesRoot.SSResponse.Children {
		cl := &c.Class
		if !isClassValid(t.logger, cl) {
			continue
		}
		cls = append(cls, cl.ID)
	}
	return cls
}

func (t *Transformer) transformEvents(eventsRoot *model.EventsRoot) ([]*pb.Event, error) {
	evs := make([]*pb.Event, 0, len(eventsRoot.SSResponse.Children))
	for _, e := range eventsRoot.SSResponse.Children {
		ev := &e.Event
		if !isEventValid(ev) {
			continue
		}
		tev, err := t.transformEvent(ev)
		if err != nil {
			return nil, err
		}
		evs = append(evs, tev)
	}
	return evs, nil
}

func (t *Transformer) transformEvent(event *model.Event) (*pb.Event, error) {
	var (
		stp = mapping.SportTypes[event.CategoryCode]
		lg  = event.TypeName
		pts = t.getParticipants(event)
	)

	sti, err := time.Parse(time.RFC3339, event.StartTime)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseTime, err)
	}

	mks, err := t.getMarkets(event)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseMarket, err)
	}

	return &pb.Event{
		// ID:           bookmaker.GenerateId(st, stp, lg, pts),
		SportType:    stp,
		Name:         event.Name,
		League:       lg,
		StartTime:    timestamppb.New(sti),
		Participants: pts,
		Markets:      mks,
	}, nil
}

func isClassValid(logger *logrus.Entry, class *model.Class) bool {
	a, err := strconv.ParseBool(class.IsActive)
	if err != nil {
		logger.WithField("class", class).Error("Failed to parse class availability")
		a = false
	}
	return a
}

func isEventValid(event *model.Event) bool {
	return event.ID != "" && event.Name != "" && len(event.StartTime) > 0
}
