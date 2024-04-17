package poller

import (
	"fmt"
	"strings"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/transform"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

const (
	classesUrl        = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/Class?translationLang=en&responseFormat=json&simpleFilter=class.isActive&simpleFilter=class.hasOpenEvent&simpleFilter=class.categoryId:equals:%v"
	classesStorageKey = "CLASSES_%s"
)

func (p *Poller) pollClasses(sportType pb.SportType) ([]byte, error) {
	rawCls, err := p.getClasses(sportType, p.config.Classes.RequestTimeout)
	if err != nil {
		return nil, err
	}
	return []byte(strings.Join(rawCls, ",")), nil
}

func (p *Poller) getClasses(sportType pb.SportType, timeout time.Duration) ([]string, error) {
	res, err := p.httpClient.Get(
		fmt.Sprintf(classesUrl, mapping.SportTypesCodes[sportType]),
		httptls.WithTimeout(timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transform.TransformClasses(res.Body)
}
