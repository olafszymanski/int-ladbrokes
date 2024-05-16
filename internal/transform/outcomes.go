package transform

import (
	"fmt"
	"strconv"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

var (
	ErrParseOutcomeAvailability   = fmt.Errorf("failed to parse outcome availability")
	ErrParsePoints                = fmt.Errorf("failed to parse points")
	ErrParsePrice                 = fmt.Errorf("failed to parse price")
	ErrParseFixedOddsAvailability = fmt.Errorf("failed to parse fixed odds availability")
)

func getOutcomes(market *model.Market, participantsOutcomeTypes map[string]pb.Outcome_OutcomeType) ([]*pb.Outcome, error) {
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

		foav, err := strconv.ParseBool(oc.FixedOddsAvail)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrParseFixedOddsAvailability, err)
		}

		tp, err := getOutcomeType(oc.Name, participantsOutcomeTypes)
		if err != nil {
			logrus.Error(err.Error())
			continue
		}

		o := &pb.Outcome{
			Type:       tp,
			ExternalId: oc.ID,
			Points:     points,
			Odds: &pb.Odds{
				Decimal:     dec,
				American:    oc.Children[0].Price.PriceAmerican,
				Numerator:   oc.Children[0].Price.PriceNum,
				Denominator: oc.Children[0].Price.PriceDen,
				IsFixed:     foav,
			},
			IsAvailable: ocav,
		}
		if o.Type == pb.Outcome_COMPETITOR {
			n := oc.Name
			o.Name = &n
		}

		outcomes = append(outcomes, o)
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

func getOutcomeType(outcomeName string, participantsOutcomeTypes map[string]pb.Outcome_OutcomeType) (pb.Outcome_OutcomeType, error) {
	var tp pb.Outcome_OutcomeType
	if t, ok := mapping.OutcomeTypes[outcomeName]; ok {
		tp = t
		return tp, nil
	} else if t, ok := participantsOutcomeTypes[outcomeName]; ok {
		tp = t
		return tp, nil
	}
	return tp, fmt.Errorf("failed to map outcome type for: %s", outcomeName)
}
