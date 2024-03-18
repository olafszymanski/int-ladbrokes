package transformer_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/transformer"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	//go:embed testdata/basketball/invalid.json
	basketballInvalidData []byte
	//go:embed testdata/basketball/empty_start_time.json
	basketballEmptyStartTimeData []byte
	//go:embed testdata/basketball/invalid_start_time.json
	basketballInvalidStartTimeData []byte
	//go:embed testdata/basketball/too_many_participants.json
	basketballTooManyParticipantsData []byte
	//go:embed testdata/basketball/invalid_points_from_market.json
	basketballInvalidPointsFromMarketData []byte
	//go:embed testdata/basketball/invalid_points_from_price.json
	basketballInvalidPointsFromPriceData []byte
	//go:embed testdata/basketball/invalid_fixed_odds_availability.json
	basketballInvalidFixedPointsAvailabilityData []byte
	//go:embed testdata/basketball/success.json
	basketballSuccessData []byte
	//go:embed testdata/basketball/success_outright.json
	basketballSuccessOutrightData []byte
)

func TestTransformEventsBasketball(t *testing.T) {
	tc := []struct {
		name        string
		data        []byte
		events      []*pb.Event
		expectedErr error
	}{
		{
			name:        "empty data",
			data:        []byte{},
			events:      nil,
			expectedErr: transformer.ErrDecodeResponse,
		},
		{
			name:        "invalid data",
			data:        basketballInvalidData,
			events:      nil,
			expectedErr: transformer.ErrDecodeResponse,
		},
		{
			name:        "empty start time",
			data:        basketballEmptyStartTimeData,
			events:      nil,
			expectedErr: transformer.ErrParseTime,
		},
		{
			name:        "invalid start time",
			data:        basketballInvalidStartTimeData,
			events:      nil,
			expectedErr: transformer.ErrParseTime,
		},
		{
			name:        "too many participants",
			data:        basketballTooManyParticipantsData,
			events:      nil,
			expectedErr: transformer.ErrTooManyParticipants,
		},
		{
			name:        "invalid points from market",
			data:        basketballInvalidPointsFromMarketData,
			events:      nil,
			expectedErr: transformer.ErrParsePoints,
		},
		{
			name:        "invalid points from price",
			data:        basketballInvalidPointsFromPriceData,
			events:      nil,
			expectedErr: transformer.ErrParsePoints,
		},
		{
			name:        "invalid fixed odds availability",
			data:        basketballInvalidFixedPointsAvailabilityData,
			events:      nil,
			expectedErr: transformer.ErrParseFixedOddsAvailability,
		},
		{
			name: "success",
			data: basketballSuccessData,
			events: []*pb.Event{
				{
					// ID:          "1",
					SportType: pb.SportType_BASKETBALL,
					Name:      "AS Monaco vs Crvena Zvezda",
					League:    "Euroleague Men",
					StartTime: timestamppb.New(time.Date(2024, 3, 7, 18, 0, 0, 0, time.UTC)),
					Participants: []*pb.Participant{
						{
							Type: getParticipantTypePtr(pb.Participant_HOME),
							Name: "AS Monaco",
						},
						{
							Type: getParticipantTypePtr(pb.Participant_AWAY),
							Name: "Crvena Zvezda",
						},
					},
					Markets: []*pb.Market{
						{
							Type: pb.MarketType_MONEYLINE.Enum(),
							Outcomes: []*pb.Outcome{
								{
									Name:   "AS Monaco",
									Points: nil,
									Odds: &pb.Odds{
										Decimal:     1.22,
										American:    "-450",
										Numerator:   "2",
										Denominator: "9",
										IsFixed:     true,
									},
									IsAvailable: true,
								},
								{
									Name:   "Crvena Zvezda",
									Points: nil,
									Odds: &pb.Odds{
										Decimal:     3.7,
										American:    "270",
										Numerator:   "27",
										Denominator: "10",
										IsFixed:     true,
									},
									IsAvailable: true,
								},
							},
						},
					},
					Link: getStringPtr("https://sports.ladbrokes.com/event/basketball/european-competitions/euroleague-men/as-monaco-v-crvena-zvezda/243810572/all-markets"),
				},
			},
			expectedErr: nil,
		},
		{
			name: "success outright",
			data: basketballSuccessOutrightData,
			events: []*pb.Event{
				{
					// ID:          "1",
					SportType: pb.SportType_BASKETBALL,
					Name:      "2023/2024 Spanish ACB",
					League:    "Spanish ACB",
					StartTime: timestamppb.New(time.Date(2024, 6, 30, 19, 0, 0, 0, time.UTC)),
					Participants: []*pb.Participant{
						{
							Type: getParticipantTypePtr(pb.Participant_COMPETITOR),
							Name: "Baskonia",
						},
						{
							Type: getParticipantTypePtr(pb.Participant_COMPETITOR),
							Name: "Basquet Girona",
						},
						{
							Type: getParticipantTypePtr(pb.Participant_COMPETITOR),
							Name: "Baxi Manresa",
						},
					},
					Markets: []*pb.Market{
						{
							Type: pb.MarketType_OUTRIGHT.Enum(),
							Outcomes: []*pb.Outcome{
								{
									Name:   "Baskonia",
									Points: nil,
									Odds: &pb.Odds{
										Decimal:     34.0,
										American:    "3300",
										Numerator:   "33",
										Denominator: "1",
										IsFixed:     true,
									},
									IsAvailable: true,
								},
								{
									Name:   "Basquet Girona",
									Points: nil,
									Odds: &pb.Odds{
										Decimal:     101.0,
										American:    "10000",
										Numerator:   "100",
										Denominator: "1",
										IsFixed:     true,
									},
									IsAvailable: true,
								},
								{
									Name:   "Baxi Manresa",
									Points: nil,
									Odds: &pb.Odds{
										Decimal:     101.0,
										American:    "10000",
										Numerator:   "100",
										Denominator: "1",
										IsFixed:     true,
									},
									IsAvailable: true,
								},
							},
						},
					},
					Link: getStringPtr("https://sports.ladbrokes.com/event/basketball/spanish/spanish-acb/2023-2024-spanish-acb/241631428/all-markets"),
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tc {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			evs, err := transformer.TransformEvents(tt.data)
			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, tt.events, evs)
		})
	}
}

func getStringPtr(s string) *string {
	return &s
}

func getParticipantTypePtr(t pb.Participant_ParticipantType) *pb.Participant_ParticipantType {
	return &t
}
