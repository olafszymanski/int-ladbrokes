package transform

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrDecodeResponse      = fmt.Errorf("decoding response failed")
	ErrParseTime           = fmt.Errorf("parsing time failed")
	ErrTooManyParticipants = fmt.Errorf("too many participants")
	ErrParseBool           = fmt.Errorf("parsing bool failed")
)

func TransformClasses(rawData []byte) ([]string, error) {
	var root model.ClassesRoot
	if err := json.Unmarshal(rawData, &root); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDecodeResponse, err)
	}
	return transformClasses(&root), nil
}

func TransformEvents(rawData []byte) ([]*pb.Event, error) {
	var root model.EventsRoot
	if err := json.Unmarshal(rawData, &root); err != nil {
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
	for _, c := range eventsRoot.SSResponse.Children {
		ev := &c.Event
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
	if len(umtps) > 0 {
		// logrus.WithField("unhandled_market_types", stringifyMarketTypes(umtps)).Warn("found unhandled market types")
	}
	return evs, nil
}

func transformEvent(event *model.Event) (*pb.Event, map[string]struct{}, error) {
	st, err := getStartTime(event.StartTime)
	if err != nil {
		return nil, nil, err
	}

	live, err := isLive(event.IsStarted)
	if err != nil {
		return nil, nil, err
	}

	pts, err := getParticipants(event)
	if err != nil {
		return nil, nil, err
	}

	mks, umtps, err := getMarkets(
		event,
		mapping.MapParticipantsToOutcomeTypes(pts),
	)
	if err != nil {
		return nil, nil, err
	}

	name, err := getName(event.Name, pts)
	if err != nil {
		return nil, nil, err
	}

	return &pb.Event{
		// ID:           bookmaker.GenerateId(st, stp, lg, pts),
		ExternalId:   event.ID,
		SportType:    mapping.SportTypes[event.CategoryCode],
		Name:         name,
		League:       event.TypeName,
		StartTime:    timestamppb.New(st),
		IsLive:       live,
		Participants: pts,
		Markets:      mks,
		Link:         getLink(event),
	}, umtps, nil
}

func isClassValid(class *model.Class) bool {
	if class.IsActive == "" {
		return false
	}

	a, err := strconv.ParseBool(class.IsActive)
	if err != nil {
		logrus.WithField("class", class).Error("failed to parse class availability")
		a = false
	}
	return a
}

func isEventValid(event *model.Event) bool {
	return event.ID != "" && event.Name != ""
}

func stringifyMarketTypes(marketTypes map[string]struct{}) string {
	var (
		i int
		s string
	)
	for k := range marketTypes {
		if i != len(marketTypes)-1 {
			s += k + ","
		}
		i++
	}
	return s
}
