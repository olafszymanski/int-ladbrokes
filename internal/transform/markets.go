package transform

import (
	"strings"

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

		m := &pb.Market{
			Type:       tp,
			ExternalId: mr.ID,
			Outcomes:   oc,
		}
		if isPlayerMarket(mr) {
			n := getPlayerName(mr.Name)
			m.Name = &n
		}

		markets = append(markets, m)
	}
	return markets, unhandledMarketTypes, nil
}

func isPlayerMarket(market *model.Market) bool {
	return market.TemplateMarketName == mapping.PlayerTotalPointsMarketType ||
		market.TemplateMarketName == mapping.PlayerTotalAssistsMarketType ||
		market.TemplateMarketName == mapping.PlayerTotalReboundsMarketType ||
		market.TemplateMarketName == mapping.PlayerTotal3PointersMarketType
}

func getPlayerName(marketName string) string {
	return strings.Split(marketName, " (")[0]
}
