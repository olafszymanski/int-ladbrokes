package client

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/transformer"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
	"github.com/sirupsen/logrus"
)

const (
	classesUrl        = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/Class?translationLang=en&responseFormat=json&simpleFilter=class.isActive&simpleFilter=class.hasOpenEvent&simpleFilter=class.categoryId:equals:%v"
	classesStorageKey = "%s_CLASSES"
	classesTTL        = int64(5 * time.Second)
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
	return transformer.TransformClasses(res.Body)
}

func (c *client) getClasses(sportType pb.SportType) (string, error) {
	cls, err := c.storage.GetAny(classesStorageKey)
	if err != nil && errors.Is(err, storage.ErrNotFound) {
		rawCls, err := c.fetchClasses(sportType)
		if err != nil {
			return "", err
		}

		cls = strings.Join(rawCls, ",")

		go func(cls any) {
			if err = c.storage.StoreAny(
				fmt.Sprintf(classesStorageKey, sportType),
				cls,
				classesTTL,
			); err != nil {
				logrus.WithError(err).Error("failed to store classes in storage")
			}
		}(cls)
	}
	return cls.(string), nil
}
