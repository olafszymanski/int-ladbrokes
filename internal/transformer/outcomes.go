package transformer

import (
	"fmt"
	"strconv"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

var (
	ErrParseOutcomeAvailability = fmt.Errorf("failed to parse outcome availability")
	ErrParsePoints              = fmt.Errorf("failed to parse points")
	ErrParsePrice               = fmt.Errorf("failed to parse price")
)

func getOutcomes(market *model.Market) ([]*pb.Outcome, error) {
	var (
		isSpread = isMarketName(market, mapping.Spread1stHalfMarketType) || isMarketName(market, mapping.Spread1stQuarterMarketType)
		points   *float64
		err      error
	)
	if !isSpread {
		points, err = getPointsFromMarket(market)
		if err != nil {
			return nil, err
		}
	}

	var outcomes []*pb.Outcome
	for _, c := range market.Children {
		oc := &c.Outcome

		dec, err := strconv.ParseFloat(oc.Children[0].Price.PriceDec, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrParsePrice, err)
		}

		if isSpread {
			points, err = getPointsFromPrice(&oc.Children[0].Price)
			if err != nil {
				return nil, err
			}
		}

		ocav, err := getOutcomeAvailability(oc)
		if err != nil {
			logrus.Error(err.Error())
		}

		outcomes = append(outcomes, &pb.Outcome{
			Name:   oc.Name,
			Points: points,
			Odds: &pb.Odds{
				Decimal:     dec,
				Numerator:   oc.Children[0].Price.PriceNum,
				Denominator: oc.Children[0].Price.PriceDen,
				American:    oc.Children[0].Price.PriceAmerican,
			},
			IsAvailable: ocav,
		})
	}
	return outcomes, nil
}

// For some reason spread points are returned with a trailing comma (everytime?) so we need to normalize them
func normalizePoints(points string) string {
	if points[len(points)-1] == ',' {
		return points[:len(points)-1]
	}
	return points
}

func getPointsFromMarket(market *model.Market) (*float64, error) {
	v := market.RawHandicapValue
	if v == "" {
		return nil, nil
	}
	p, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, fmt.Errorf("%w from market: %s", ErrParsePoints, err)
	}
	return &p, nil
}

func getPointsFromPrice(price *model.Price) (*float64, error) {
	v := price.HandicapValueDec
	if v == "" {
		v = price.RawHandicapValue
		if v == "" {
			return nil, nil
		}
	}
	v = normalizePoints(v)

	p, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, fmt.Errorf("%w from price: %s", ErrParsePoints, err)
	}
	return &p, nil
}

func getOutcomeAvailability(outcome *model.Outcome) (bool, error) {
	av, err := strconv.ParseBool(outcome.IsAvailable)
	if err == nil {
		return av, nil
	}

	dp, err := strconv.ParseBool(outcome.IsDisplayed)
	if err != nil {
		return false, fmt.Errorf("%w: %s", ErrParseOutcomeAvailability, err)
	}
	return dp && outcome.OutcomeStatusCode != model.SuspendedOutcomeCode, nil
}
