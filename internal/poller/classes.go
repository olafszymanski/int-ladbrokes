package poller

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/transform"
	sdkHttp "github.com/olafszymanski/int-sdk/http"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

const (
	classesUrl        = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/Class?translationLang=en&responseFormat=json&simpleFilter=class.isActive&simpleFilter=class.hasOpenEvent&simpleFilter=class.categoryId:equals:%v"
	classesStorageKey = "CLASSES_%s"
)

func (p *Poller) pollClasses(ctx context.Context, logger *logrus.Entry, sportType pb.SportType) error {
	var (
		startTime time.Time
		classesCh = make(chan []byte)
	)
	defer close(classesCh)

	for {
		startTime = time.Now()
		go func() {
			cls, err := p.fetchClasses(sportType)
			if err != nil {
				logger.WithError(err).Error("polling classes failed")
				return
			}
			if len(cls) == 0 {
				logger.Warn("no classes polled")
				return
			}
			classesCh <- cls
		}()

		select {
		case cls := <-classesCh:
			logger.WithField("classes_length", len(cls)).Debug("classes polled")
			if err := p.storage.StoreClasses(ctx, fmt.Sprintf(classesStorageKey, sportType), cls); err != nil {
				return err
			}
			<-time.After(p.config.Classes.RequestInterval - time.Since(startTime))
		case <-time.After(p.config.Classes.RequestInterval):
			logger.Warn("classes polling took longer than expected")
		}
	}
}

func (p *Poller) fetchClasses(sportType pb.SportType) ([]byte, error) {
	rawCls, err := p.getClasses(sportType, p.config.Classes.RequestTimeout)
	if err != nil {
		return nil, err
	}
	return []byte(strings.Join(rawCls, ",")), nil
}

func (p *Poller) getClasses(sportType pb.SportType, timeout time.Duration) ([]string, error) {
	res, err := p.httpClient.Do(&sdkHttp.Request{
		Method:  http.MethodGet,
		URL:     fmt.Sprintf(classesUrl, mapping.SportTypesCodes[sportType]),
		Timeout: timeout,
	})
	if err != nil {
		return nil, err
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transform.TransformClasses(res.Body)
}
