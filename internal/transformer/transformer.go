package transformer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

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
		logrus.WithField("unhandled_market_types", stringifyMarketTypes(umtps)).Warn("found unhandled market types")
	}
	return evs, nil
}

func transformEvent(event *model.Event) (*pb.Event, map[string]struct{}, error) {
	var (
		eId = event.ID
		stp = mapping.SportTypes[event.CategoryCode]
		lg  = event.TypeName
		pts = getParticipants(event)
	)

	if len(strings.TrimSpace(event.StartTime)) == 0 {
		return nil, nil, fmt.Errorf("%w: empty start time", ErrParseTime)
	}
	sti, err := time.Parse(time.RFC3339, event.StartTime)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrParseTime, err)
	}

	otm := mapping.MapParticipantsToOutcomeTypes(pts)
	mks, umtps, err := getMarkets(event, otm)
	if err != nil {
		return nil, nil, err
	}

	name := event.Name
	if len(pts) > 0 && *pts[0].Type != pb.Participant_COMPETITOR {
		if len(pts) > 2 {
			return nil, nil, fmt.Errorf("%w: expected 2", ErrTooManyParticipants)
		}
		name = fmt.Sprintf("%s vs %s", pts[0].Name, pts[1].Name)
	}

	l := getEventLink(event)

	return &pb.Event{
		// ID:           bookmaker.GenerateId(st, stp, lg, pts),
		ExternalId:   &eId,
		SportType:    stp,
		Name:         name,
		League:       lg,
		StartTime:    timestamppb.New(sti),
		Participants: pts,
		Markets:      mks,
		Link:         &l,
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
	var s string
	for k := range marketTypes {
		s += k + ","
	}
	return s
}

func getEventLink(event *model.Event) string {
	return fmt.Sprintf(
		"https://sports.ladbrokes.com/event/%s/%s/%s/%s/%s/all-markets",
		normalizeLinkPart(event.CategoryName),
		normalizeLinkPart(event.ClassName),
		normalizeLinkPart(event.TypeName),
		normalizeLinkPart(event.Name),
		event.ID,
	)
}

func normalizeLinkPart(part string) string {
	p := strings.ReplaceAll(part, " ", "-")
	p = strings.ReplaceAll(p, "/", "-")
	p = strings.ToLower(p)
	return p
}
