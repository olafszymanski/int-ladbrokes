package mapping

type UpdateType uint8

const (
	UnknownUpdateType UpdateType = iota
	EventUpdateType
	MarketUpdateType
	SelectionUpdateType
	PriceUpdateType
)

var UpdateTypes = map[string]UpdateType{
	"EVENT": MarketUpdateType,
	"EVMKT": MarketUpdateType,
	"SELCN": SelectionUpdateType,
	"PRICE": PriceUpdateType,
}
