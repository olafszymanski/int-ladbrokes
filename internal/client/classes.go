package client

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/transform"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
	"github.com/sirupsen/logrus"
)

const (
	classesUrl        = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/Class?translationLang=en&responseFormat=json&simpleFilter=class.isActive&simpleFilter=class.hasOpenEvent&simpleFilter=class.categoryId:equals:%v"
	classesStorageKey = "%s_CLASSES"
	classesTTL        = 5 * time.Second
)

func (c *client) fetchClasses(sportType pb.SportType) ([]string, error) {
	res, err := c.httpClient.Get(
		fmt.Sprintf(classesUrl, mapping.SportTypesCodes[sportType]),
		httptls.WithTimeout(2000),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transform.TransformClasses(res.Body)
}

func (c *client) getClasses(ctx context.Context, sportType pb.SportType) (string, error) {
	cls, err := c.storage.Get(ctx, classesStorageKey)
	if err != nil && errors.Is(err, storage.ErrNotFound) {
		rawCls, err := c.fetchClasses(sportType)
		if err != nil {
			return "", err
		}

		cls = []byte(strings.Join(rawCls, ","))

		go func() {
			if err = c.storage.Store(
				ctx,
				fmt.Sprintf(classesStorageKey, sportType),
				cls,
				classesTTL,
			); err != nil {
				logrus.WithError(err).Error("failed to store classes in storage")
			}
		}()
	}
	return string(cls), nil
}
