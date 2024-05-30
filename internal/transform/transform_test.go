package transform_test

import (
	_ "embed"
	"testing"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/transform"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	//go:embed testdata/events/invalid.json
	eventsInvalidData []byte
	//go:embed testdata/events/empty_start_time.json
	eventsEmptyStartTimeData []byte
	//go:embed testdata/events/invalid_start_time.json
	eventsInvalidStartTimeData []byte
	//go:embed testdata/events/too_many_participants.json
	eventsTooManyParticipantsData []byte
	//go:embed testdata/events/invalid_points_from_market.json
	eventsInvalidPointsFromMarketData []byte
	//go:embed testdata/events/invalid_points_from_price.json
	eventsInvalidPointsFromPriceData []byte
	//go:embed testdata/events/invalid_fixed_odds_availability.json
	eventsInvalidFixedPointsAvailabilityData []byte
	//go:embed testdata/events/success.json
	eventsSuccessData []byte
	//go:embed testdata/events/success_outright.json
	eventsSuccessOutrightData []byte

	//go:embed testdata/updates/success.json
	updatesSuccessData []byte
)

// TODO: Add test for classes and updates

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
			expectedErr: transform.ErrDecodeResponse,
		},
		{
			name:        "invalid data",
			data:        eventsInvalidData,
			events:      nil,
			expectedErr: transform.ErrDecodeResponse,
		},
		{
			name:        "empty start time",
			data:        eventsEmptyStartTimeData,
			events:      nil,
			expectedErr: transform.ErrParseTime,
		},
		{
			name:        "invalid start time",
			data:        eventsInvalidStartTimeData,
			events:      nil,
			expectedErr: transform.ErrParseTime,
		},
		{
			name:        "too many participants",
			data:        eventsTooManyParticipantsData,
			events:      nil,
			expectedErr: transform.ErrTooManyParticipants,
		},
		{
			name:        "invalid points from market",
			data:        eventsInvalidPointsFromMarketData,
			events:      nil,
			expectedErr: transform.ErrParsePoints,
		},
		{
			name:        "invalid points from price",
			data:        eventsInvalidPointsFromPriceData,
			events:      nil,
			expectedErr: transform.ErrParsePoints,
		},
		{
			name:        "invalid fixed odds availability",
			data:        eventsInvalidFixedPointsAvailabilityData,
			events:      nil,
			expectedErr: transform.ErrParseFixedOddsAvailability,
		},
		{
			name: "success",
			data: eventsSuccessData,
			events: []*pb.Event{
				{
					// ID:          "1",
					ExternalId: "243810572",
					SportType:  pb.SportType_BASKETBALL,
					Name:       "AS Monaco vs Crvena Zvezda",
					League:     "Euroleague Men",
					StartTime:  timestamppb.New(time.Date(2024, 3, 7, 18, 0, 0, 0, time.UTC)),
					IsLive:     true,
					Participants: []*pb.Participant{
						{
							Type: pb.Participant_HOME,
							Name: "AS Monaco",
						},
						{
							Type: pb.Participant_AWAY,
							Name: "Crvena Zvezda",
						},
					},
					Markets: []*pb.Market{
						{
							Type: pb.MarketType_MONEYLINE,
							Name: nil,
							Outcomes: []*pb.Outcome{
								{
									Type:   pb.Outcome_HOME,
									Name:   nil,
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
									Type:   pb.Outcome_AWAY,
									Name:   nil,
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
						{
							Type: pb.MarketType_PLAYER_TOTAL_POINTS,
							Name: getStringPtr("Donatas Motiejunas"),
							Outcomes: []*pb.Outcome{
								{
									Type:   pb.Outcome_OVER,
									Name:   nil,
									Points: getFloat64Ptr(7.5),
									Odds: &pb.Odds{
										Decimal:     1.83,
										American:    "-121",
										Numerator:   "83",
										Denominator: "100",
										IsFixed:     true,
									},
									IsAvailable: true,
								},
								{
									Type:   pb.Outcome_UNDER,
									Name:   nil,
									Points: getFloat64Ptr(7.5),
									Odds: &pb.Odds{
										Decimal:     1.83,
										American:    "-121",
										Numerator:   "83",
										Denominator: "100",
										IsFixed:     true,
									},
									IsAvailable: true,
								},
							},
						},
					},
					Link: "https://sports.ladbrokes.com/event/basketball/european-competitions/euroleague-men/as-monaco-v-crvena-zvezda/243810572/all-markets",
				},
			},
			expectedErr: nil,
		},
		{
			name: "success outright",
			data: eventsSuccessOutrightData,
			events: []*pb.Event{
				{
					// ID:          "1",
					ExternalId: "241631428",
					SportType:  pb.SportType_BASKETBALL,
					Name:       "2023/2024 Spanish ACB",
					League:     "Spanish ACB",
					StartTime:  timestamppb.New(time.Date(2024, 6, 30, 19, 0, 0, 0, time.UTC)),
					Participants: []*pb.Participant{
						{
							Type: pb.Participant_COMPETITOR,
							Name: "Baskonia",
						},
						{
							Type: pb.Participant_COMPETITOR,
							Name: "Basquet Girona",
						},
						{
							Type: pb.Participant_COMPETITOR,
							Name: "Baxi Manresa",
						},
					},
					Markets: []*pb.Market{
						{
							Type: pb.MarketType_OUTRIGHT,
							Name: nil,
							Outcomes: []*pb.Outcome{
								{
									Type:   pb.Outcome_COMPETITOR,
									Name:   getStringPtr("Baskonia"),
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
									Type:   pb.Outcome_COMPETITOR,
									Name:   getStringPtr("Basquet Girona"),
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
									Type:   pb.Outcome_COMPETITOR,
									Name:   getStringPtr("Baxi Manresa"),
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
					Link: "https://sports.ladbrokes.com/event/basketball/spanish/spanish-acb/2023-2024-spanish-acb/241631428/all-markets",
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tc {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			evs, err := transform.TransformEvents(tt.data)
			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, tt.events, evs)
		})
	}
}

func getStringPtr(s string) *string {
	return &s
}

func getFloat64Ptr(f float64) *float64 {
	return &f
}

func TestTransformUpdates(t *testing.T) {
	tc := []struct {
		name        string
		data        []byte
		update      *transform.Update
		expectedErr error
	}{
		{
			name: "success",
			data: updatesSuccessData,
			update: &transform.Update{
				Data: map[mapping.UpdateType][]*transform.UpdateData{
					mapping.PriceUpdateType: {
						{
							ID:      "1",
							RawData: []byte(`MSEVENT0244772570!!!!'o0Lm8GsPRICE237261105000001e00001e{"lp_num": "4", "lp_den": "7"}`),
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tc {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			u, err := transform.TransformUpdates(tt.data)
			require.ErrorIs(t, err, tt.expectedErr)
			require.Equal(t, tt.update, u)
		})
	}
}
