package transform

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/model"
)

const (
	majorEventSeparator   = "MSEVENT0" // TODO: Change names
	minorEventSeparator   = "MsEVENT0" // TODO: Fix id starting with 0
	idLength              = 10
	requestBodyPartLength = 6
	updateTypeLength      = 5

	/*
		Example response: MSEVENT0244772570!!!!'o0:FCGsPRICE237261429800001e00001e{"lp_num": "3", "lp_den": "4"}

		After splitting with splitUpdates() we are left with:
		0244772570!!!!'o0:FCGsPRICE237261429800001e00001e{"lp_num": "3", "lp_den": "4"}

		0244772570 !!!! 'o0:FC Gs PRICE 2372614298 00001e00001e {"lp_num": "3", "lp_den": "4"}

		0244772570          = id of the event (has length o 9, we remove 0 from the beginning)
		!!!!                = we don't need this part
		'o0:FC              = a part of next request body
		Gs		           = we don't need this part
		PRICE	           = update type
		2372614298          = id of the update (has length of 10)
		00001e00001e        = we don't need this part (has length of 12)
		{"lp_num": "3",...} = actual update data
	*/

	updateTypeStartIndex       = idLength + 4 + requestBodyPartLength + 2 - 1
	updateIdStartIndex         = idLength + 4 + requestBodyPartLength + 2 + updateTypeLength - 1
	updateDataStartIndex       = idLength + 4 + requestBodyPartLength + 2 + updateTypeLength + idLength + 12 - 1
	requestBodyPartsStartIndex = idLength + 4 - 1

	suspendedStatus = "S"
	notDisplayed    = "N"
)

var defaultLastRequestBodyPart = []byte("!!!!!0")

type updateType interface {
	model.EventUpdate | model.MarketUpdate | model.PriceUpdate | model.SelectionUpdate
}

func UnmarshalUpdate[T updateType](rawData []byte) (*T, error) {
	t := new(T)
	if err := json.Unmarshal(rawData, t); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDecodeResponse, err)
	}
	return t, nil
}

func IsEventFinished(update *model.EventUpdate) bool {
	return update.Status == suspendedStatus && update.Displayed == notDisplayed
}

func IsMarketSuspended(update *model.MarketUpdate) bool {
	return update.Status == suspendedStatus
}

func IsMarketRemoved(update *model.MarketUpdate) bool {
	return update.Status == suspendedStatus && update.Displayed == notDisplayed
}

func IsSelectionSuspended(update *model.SelectionUpdate) bool {
	return update.Status == suspendedStatus
}

// splits rawData into individual raw updates for further processing
func splitRawData(rawData []byte) [][]byte {
	ups := bytes.Split(rawData, []byte(majorEventSeparator))
	updates := make([][]byte, 0, len(ups))
	for _, u := range ups {
		nu := bytes.Split(u, []byte(minorEventSeparator))
		updates = append(updates, nu...)
	}
	return updates[1:] // 0 element is empty
}

func getUpdateType(rawData []byte) (mapping.UpdateType, error) {
	raw := string(rawData[updateTypeStartIndex : updateTypeStartIndex+updateTypeLength])
	ut, ok := mapping.UpdateTypes[string(raw)]
	if !ok {
		return mapping.UnknownUpdateType, fmt.Errorf("update type '%s' not found", raw)
	}
	return ut, nil
}

// TODO: FIX
func getEventId(rawData []byte) string {
	return string(rawData[:idLength-1])
}

func getUpdateId(rawData []byte, updateType mapping.UpdateType) string {
	// EventUpdateType and MarketUpdateType updates have the same ID length of 9, e.g: 02372538782, thus we have to remove 0
	if updateType == mapping.EventUpdateType || updateType == mapping.MarketUpdateType {
		return string(rawData[updateIdStartIndex+1 : updateIdStartIndex+idLength])
	}
	return string(rawData[updateIdStartIndex : updateIdStartIndex+idLength])
}

func getUpdateDataId(rawData []byte, updateType mapping.UpdateType) (string, error) {
	switch updateType {
	case mapping.EventUpdateType:
		return getEventId(rawData), nil
	case mapping.SelectionUpdateType, mapping.PriceUpdateType, mapping.MarketUpdateType:
		return getUpdateId(rawData, updateType), nil
	default:
		return "", fmt.Errorf("unknown update type")
	}
}

func getUpdateData(rawData []byte) []byte {
	return rawData[updateDataStartIndex:]
}

func getRequestBodyParts(rawUpdates [][]byte) []string {
	var (
		first = rawUpdates[0][requestBodyPartsStartIndex : requestBodyPartsStartIndex+requestBodyPartLength]
		last  = defaultLastRequestBodyPart
	)
	if len(rawUpdates) != 1 {
		last = rawUpdates[len(rawUpdates)-1][requestBodyPartsStartIndex : requestBodyPartsStartIndex+requestBodyPartLength]
	}
	return []string{string(first), string(last)}
}
