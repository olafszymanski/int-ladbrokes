package mapping

import (
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

var OutcomeTypes = map[string]pb.Outcome_OutcomeType{
	"Draw":  pb.Outcome_DRAW,
	"Over":  pb.Outcome_OVER,
	"Under": pb.Outcome_UNDER,
	"Odd":   pb.Outcome_ODD,
	"Even":  pb.Outcome_EVEN,
}

func MapParticipantsToOutcomeTypes(participants []*pb.Participant) map[string]pb.Outcome_OutcomeType {
	pts := make(map[string]pb.Outcome_OutcomeType, len(participants))
	for _, p := range participants {
		var t pb.Outcome_OutcomeType
		switch p.Type {
		case pb.Participant_HOME:
			t = pb.Outcome_HOME
		case pb.Participant_AWAY:
			t = pb.Outcome_AWAY
		case pb.Participant_COMPETITOR:
			t = pb.Outcome_COMPETITOR
		default:
			logrus.WithField("participant", p).Error("failed to map outcome type, unknown participant type")
		}
		pts[p.Name] = t
	}
	return pts
}
