package transformer

import (
	"fmt"
	"strconv"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

var ErrParseMarketAvailability = fmt.Errorf("failed to parse market availability")

func (t *Transformer) getMarkets(event *model.Event) ([]*pb.Market, error) {
	var markets []*pb.Market
	for _, m := range event.Children {
		mr := &m.Market
		t.logger = t.logger.WithField("market", mr)

		tp, ok := mapping.MarketTypes[mr.TemplateMarketName]
		if !ok {
			t.logger.Warn("Found unhandled market type")
			continue
		}

		oc, err := t.getOutcomes(mr)
		if err != nil {
			return nil, err
		}

		av, err := getMarketAvailability(mr)
		if err != nil {
			t.logger.Warn(err.Error())
		}

		markets = append(markets, &pb.Market{
			Type:        tp,
			Outcomes:    oc,
			IsAvailable: av,
		})
	}
	return markets, nil
}

func getMarketAvailability(market *model.Market) (bool, error) {
	av, err := strconv.ParseBool(market.IsAvailable)
	if err == nil {
		return av, nil
	}

	lpav, _ := strconv.ParseBool(market.IsLpAvailable)
	dp, _ := strconv.ParseBool(market.IsDisplayed)
	bt, _ := strconv.ParseBool(market.IsBettable)
	if lpav && dp && bt {
		return true, nil
	}
	return false, ErrParseMarketAvailability
}
