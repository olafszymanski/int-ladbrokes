package transform

import (
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

func getParticipants(event *model.Event) ([]*pb.Participant, error) {
	pts := make([]*pb.Participant, 0)
	for _, ec := range event.Children {
		mr := &ec.Market
		if isMarketName(mr, mapping.MoneyLineMarketType) ||
			isMarketName(mr, mapping.OutrightMarketType) {
			for _, mc := range mr.Children {
				oc := &mc.Outcome
				t := getParticipantType(oc.OutcomeMeaningMinorCode)
				pts = append(pts, &pb.Participant{
					Type: &t,
					Name: oc.Name,
				})
			}
		}
	}
	if len(pts) > 2 && *pts[0].Type != pb.Participant_COMPETITOR {
		return nil, fmt.Errorf("%w: expected 2", ErrTooManyParticipants)
	}
	return pts, nil
}

func isMarketName(market *model.Market, name string) bool {
	return market.TemplateMarketName == name
}

func getParticipantType(code string) pb.Participant_ParticipantType {
	switch code {
	case model.HomeOutcomeCode:
		return pb.Participant_HOME
	case model.AwayOutcomeCode:
		return pb.Participant_AWAY
	default:
		return pb.Participant_COMPETITOR
	}
}
