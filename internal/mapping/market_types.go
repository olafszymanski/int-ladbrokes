package mapping

import "github.com/olafszymanski/int-sdk/integration/pb"

const (
	MatchBettingMarketType = "Match Betting"

	MoneyLineMarketType           = "Money Line"
	Moneyline1stQuarterMarketType = "Quarter Money Line"
	Moneyline1stHalfMarketType    = "Half Money Line"

	RaceToXPointsMarketType = "Race to X points"

	OutrightMarketType = "Outright"

	TotalPointsMarketType                   = "Total Points"
	TotalPointsOddEvenMarketType            = "Total points - odd/even"
	TotalPoints1stQuarterMarketType         = "Quarter Total Points"
	TotalPoints1stQuarterOddEvenMarketType  = "Quarter Total Points Odd/Even"
	TotalPoints1stHalfMarketType            = "Half Total Points"
	TotalPoints1stHalfOddEvenMarketType     = "Half Total Points Odd/Even"
	TotalPointsHomeTeamMarketType           = "Home team total points"
	TotalPointsHomeTeam1stQuarterMarketType = "Quarter Home Team Total Points"
	TotalPointsHomeTeam1stHalfMarketType    = "Half Home Team Total Points"
	TotalPointsAwayTeamMarketType           = "Away team total points"
	TotalPointsAwayTeam1stQuarterMarketType = "Quarter Away Team Total Points"
	TotalPointsAwayTeam1stHalfMarketType    = "Half Away Team Total Points"

	Spread1stQuarterMarketType = "Quarter Spread"
	Spread1stHalfMarketType    = "Half Spread"

	PlayerTotalPointsMarketType    = "Player - Total Points"
	PlayerTotalAssistsMarketType   = "Player - Total Assists"
	PlayerTotalReboundsMarketType  = "Player - Total Rebounds"
	PlayerTotal3PointersMarketType = "Player - Total 3 Pointers"
)

var MarketTypes = map[string]pb.MarketType{
	MatchBettingMarketType: pb.MarketType_MATCH_BETTING,

	MoneyLineMarketType:           pb.MarketType_MONEYLINE,
	Moneyline1stQuarterMarketType: pb.MarketType_MONEYLINE_1ST_QUARTER,
	Moneyline1stHalfMarketType:    pb.MarketType_MONEYLINE_1ST_HALF,

	RaceToXPointsMarketType: pb.MarketType_RACE_TO_X_POINTS,

	OutrightMarketType: pb.MarketType_OUTRIGHT,

	TotalPointsMarketType:                   pb.MarketType_TOTAL_POINTS,
	TotalPointsOddEvenMarketType:            pb.MarketType_TOTAL_POINTS_ODD_EVEN,
	TotalPoints1stQuarterMarketType:         pb.MarketType_TOTAL_POINTS_1ST_QUARTER,
	TotalPoints1stQuarterOddEvenMarketType:  pb.MarketType_TOTAL_POINTS_1ST_QUARTER_ODD_EVEN,
	TotalPoints1stHalfMarketType:            pb.MarketType_TOTAL_POINTS_1ST_HALF,
	TotalPoints1stHalfOddEvenMarketType:     pb.MarketType_TOTAL_POINTS_1ST_HALF_ODD_EVEN,
	TotalPointsHomeTeamMarketType:           pb.MarketType_TOTAL_POINTS_HOME_TEAM,
	TotalPointsHomeTeam1stQuarterMarketType: pb.MarketType_TOTAL_POINTS_HOME_TEAM_1ST_QUARTER,
	TotalPointsHomeTeam1stHalfMarketType:    pb.MarketType_TOTAL_POINTS_HOME_TEAM_1ST_HALF,
	TotalPointsAwayTeamMarketType:           pb.MarketType_TOTAL_POINTS_AWAY_TEAM,
	TotalPointsAwayTeam1stQuarterMarketType: pb.MarketType_TOTAL_POINTS_AWAY_TEAM_1ST_QUARTER,
	TotalPointsAwayTeam1stHalfMarketType:    pb.MarketType_TOTAL_POINTS_AWAY_TEAM_1ST_HALF,

	Spread1stQuarterMarketType: pb.MarketType_SPREAD_1ST_QUARTER,
	Spread1stHalfMarketType:    pb.MarketType_SPREAD_1ST_HALF,

	PlayerTotalPointsMarketType:    pb.MarketType_PLAYER_TOTAL_POINTS,
	PlayerTotalAssistsMarketType:   pb.MarketType_PLAYER_TOTAL_ASSISTS,
	PlayerTotalReboundsMarketType:  pb.MarketType_PLAYER_TOTAL_REBOUNDS,
	PlayerTotal3PointersMarketType: pb.MarketType_PLAYER_TOTAL_3_POINTERS,
}
