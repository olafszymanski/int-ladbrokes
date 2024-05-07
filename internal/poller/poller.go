package poller

import (
	"context"
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
	"github.com/sirupsen/logrus"
)

var (
	ErrRequest              = fmt.Errorf("request failed")
	ErrUnexpectedStatusCode = fmt.Errorf("unexpected status code")
)

type Poller struct {
	config     *config.Config
	httpClient *httptls.HTTPClient
	storage    storage.Storager
	errCh      chan error
}

func NewPoller(config *config.Config, httpClient *httptls.HTTPClient, storage storage.Storager) (*Poller, error) {
	return &Poller{
		config:     config,
		httpClient: httpClient,
		storage:    storage,
		errCh:      make(chan error),
	}, nil
}

func (p *Poller) Run(ctx context.Context, sportType pb.SportType) error {
	defer close(p.errCh)

	logger := logrus.WithField("sport_type", sportType)

	go p.pollClasses(ctx, logger, sportType)
	go p.pollEvents(ctx, logger, sportType)

	if err := <-p.errCh; err != nil {
		return err
	}
	return nil
}
