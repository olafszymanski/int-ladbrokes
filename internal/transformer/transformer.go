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

func TransformClasses(reader io.Reader) ([]string, error) {
	var root model.ClassesRoot
	if err := json.NewDecoder(reader).Decode(&root); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDecodeResponse, err)
	}
	return transformClasses(&root), nil
}

func TransformEvents(reader io.Reader) ([]*pb.Event, error) {
	var root model.EventsRoot
	if err := json.NewDecoder(reader).Decode(&root); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDecodeResponse, err)
	}
	return transformEvents(&root)
}

func transformClasses(classesRoot *model.ClassesRoot) []string {
	cls := make([]string, 0, len(classesRoot.SSResponse.Children))
	for _, c := range classesRoot.SSResponse.Children {
		cl := &c.Class
		if !isClassValid(cl) {
			continue
		}
		cls = append(cls, cl.ID)
	}
	return cls
}

func transformEvents(eventsRoot *model.EventsRoot) ([]*pb.Event, error) {
	var (
		evs   = make([]*pb.Event, 0, len(eventsRoot.SSResponse.Children))
		umtps = make(map[string]struct{})
	)
	for _, e := range eventsRoot.SSResponse.Children {
		ev := &e.Event
		if !isEventValid(ev) {
			continue
		}
		tev, u, err := transformEvent(ev)
		if err != nil {
			return nil, err
		}
		evs = append(evs, tev)
		for k := range u {
			umtps[k] = struct{}{}
		}
	}
	logrus.WithField("unhandled_market_types", stringifyMarketTypes(umtps)).Warn("Found unhandled market types")
	return evs, nil
}

func transformEvent(event *model.Event) (*pb.Event, map[string]struct{}, error) {
	var (
		stp = mapping.SportTypes[event.CategoryCode]
		lg  = event.TypeName
		pts = getParticipants(event)
	)

	sti, err := time.Parse(time.RFC3339, event.StartTime)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrParseTime, err)
	}

	mks, umtps, err := getMarkets(event)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrParseMarket, err)
	}

	return &pb.Event{
		// ID:           bookmaker.GenerateId(st, stp, lg, pts),
		SportType:    stp,
		Name:         event.Name,
		League:       lg,
		StartTime:    timestamppb.New(sti),
		Participants: pts,
		Markets:      mks,
	}, umtps, nil
}

func isClassValid(class *model.Class) bool {
	a, err := strconv.ParseBool(class.IsActive)
	if err != nil {
		logrus.WithField("class", class).Error("Failed to parse class availability")
		a = false
	}
	return a
}

func isEventValid(event *model.Event) bool {
	return event.ID != "" && event.Name != "" && len(event.StartTime) > 0
}

func stringifyMarketTypes(marketTypes map[string]struct{}) string {
	var s string
	for k := range marketTypes {
		s += k + ","
	}
	return s
}
