package transform

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

func getStartTime(startTime string) (time.Time, error) {
	if len(strings.TrimSpace(startTime)) == 0 {
		return time.Time{}, fmt.Errorf("%w: empty start time", ErrParseTime)
	}
	sti, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrParseTime, err)
	}
	return sti, nil
}

func getName(name string, participants []*pb.Participant) (string, error) {
	if len(participants) > 0 && participants[0].Type != pb.Participant_COMPETITOR {
		name = fmt.Sprintf("%s vs %s", participants[0].Name, participants[1].Name)
	}
	return name, nil
}

func isLive(isStarted string) (bool, error) {
	if isStarted == "" {
		return false, nil
	}
	l, err := strconv.ParseBool(isStarted)
	if err != nil {
		return false, fmt.Errorf("%w: %s", ErrParseBool, err)
	}
	return l, nil
}

func getLink(event *model.Event) string {
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
