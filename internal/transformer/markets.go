package transformer

import (
	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

// getMarkets returns a slice of markets for the given event, unhandled market types and an error if any
func getMarkets(event *model.Event, participantsOutcomeTypes map[string]pb.Outcome_OutcomeType) ([]*pb.Market, map[string]struct{}, error) {
	var (
		markets              []*pb.Market
		unhandledMarketTypes = map[string]struct{}{}
	)
	for _, c := range event.Children {
		mr := &c.Market

		tp, ok := mapping.MarketTypes[mr.TemplateMarketName]
		if !ok {
			unhandledMarketTypes[mr.TemplateMarketName] = struct{}{}
			continue
		}

		oc, err := getOutcomes(mr, participantsOutcomeTypes)
		if err != nil {
			return nil, nil, err
		}

		markets = append(markets, &pb.Market{
			Type:     &tp,
			Outcomes: oc,
		})
	}
	return markets, unhandledMarketTypes, nil
}
