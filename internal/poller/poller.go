package poller

import (
	"context"
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-ladbrokes/internal/storage"
	"github.com/olafszymanski/int-sdk/http"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/sirupsen/logrus"
)

var ErrUnexpectedStatusCode = fmt.Errorf("unexpected status code")

type Poller struct {
	config     *config.Config
	httpClient http.Doer
	storage    *storage.Storage
}

func NewPoller(config *config.Config, httpClient http.Doer, storage *storage.Storage) (*Poller, error) {
	return &Poller{
		config:     config,
		httpClient: httpClient,
		storage:    storage,
	}, nil
}

func (p *Poller) Run(ctx context.Context, sportType pb.SportType) error {
	var (
		logger = logrus.WithField("sport_type", sportType)
		errCh  = make(chan error)
	)
	defer close(errCh)

	go func() {
		if err := p.pollClasses(ctx, logger, sportType); err != nil {
			errCh <- err
			return
		}
	}()
	go func() {
		if err := p.pollLiveEvents(ctx, logger, sportType); err != nil {
			errCh <- err
			return
		}
	}()
	go func() {
		if err := p.pollPreMatchEvents(ctx, logger, sportType); err != nil {
			errCh <- err
			return
		}
	}()
	go func() {
		if err := p.pollUpdates(ctx, logger, sportType); err != nil {
			errCh <- err
			return
		}
	}()

	if err := <-errCh; err != nil {
		return err
	}
	return nil
}
