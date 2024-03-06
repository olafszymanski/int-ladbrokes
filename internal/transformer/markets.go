package transformer

import (
	"fmt"
	"strconv"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

var ErrParseMarketAvailability = fmt.Errorf("failed to parse market availability")

// getMarkets returns a slice of markets for the given event, unhandled market types and an error if any
func getMarkets(event *model.Event) ([]*pb.Market, map[string]struct{}, error) {
	var (
		markets              []*pb.Market
		unhandledMarketTypes = map[string]struct{}{}
	)
	for _, c := range event.Children {
		mr := &c.Market
		logger := logrus.WithField("market", mr)

		tp, ok := mapping.MarketTypes[mr.TemplateMarketName]
		if !ok {
			unhandledMarketTypes[mr.TemplateMarketName] = struct{}{}
			continue
		}

		oc, err := getOutcomes(mr)
		if err != nil {
			return nil, nil, err
		}

		av, err := getMarketAvailability(mr)
		if err != nil {
			logger.Warn(err.Error())
		}

		markets = append(markets, &pb.Market{
			Type:        tp,
			Outcomes:    oc,
			IsAvailable: av,
		})
	}
	return markets, unhandledMarketTypes, nil
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
