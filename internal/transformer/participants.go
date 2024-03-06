package transformer

import (
	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

func (t *Transformer) getParticipants(event *model.Event) []*pb.Participant {
	pts := make([]*pb.Participant, 0)
	for _, m := range event.Children {
		mr := &m.Market
		if isMarketName(mr, mapping.MoneyLineMarketType) ||
			isMarketName(mr, mapping.OutrightMarketType) {
			for _, c := range mr.Children {
				oc := &c.Outcome
				pts = append(pts, &pb.Participant{
					Type: getParticipantType(oc.OutcomeMeaningMinorCode),
					Name: oc.Name,
				})
			}
		}
	}
	return pts
}

func isMarketName(market *model.Market, name string) bool {
	return market.TemplateMarketName == name
}

func getParticipantType(code string) pb.ParticipantType {
	switch code {
	case model.HomeOutcomeCode:
		return pb.ParticipantType_HOME
	case model.AwayOutcomeCode:
		return pb.ParticipantType_AWAY
	default:
		return pb.ParticipantType_COMPETITOR
	}
}
