package mapping

import "github.com/olafszymanski/int-sdk/integration/pb"

var SportTypes = map[string]pb.SportType{
	"BASKETBALL": pb.SportType_BASKETBALL,
}

var SportTypesCodes = map[pb.SportType]int{
	pb.SportType_BASKETBALL: 6,
}
