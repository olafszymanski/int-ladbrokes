package client

import (
	"context"
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
)

var (
	ErrRequest              = fmt.Errorf("request failed")
	ErrUnexpectedStatusCode = fmt.Errorf("unexpected status code")
)

type client struct {
	httpClient *httptls.HTTPClient
	config     *config.Config
	storage    storage.Storager
	pb.UnimplementedIntegrationServer
}

func NewClient(cfg *config.Config, storage storage.Storager) pb.IntegrationServer {
	return &client{
		httpClient: httptls.NewHTTPClient(),
		config:     cfg,
		storage:    storage,
	}
}

func (c *client) GetLive(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	cls, err := c.getClasses(request.SportType)
	if err != nil {
		return nil, err
	}

	evs, err := c.getEvents(liveEventsUrl, cls, c.config.MaxLiveEventsConcurrentRequests, c.config.LiveEventsRequestTimeout)
	return &pb.Response{
		Events: evs,
	}, err
}

func (c *client) GetPreMatch(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	cls, err := c.getClasses(request.SportType)
	if err != nil {
		return nil, err
	}

	evs, err := c.getEvents(preMatchEventsUrl, cls, c.config.MaxPreMatchEventsConcurrentRequests, c.config.PreMatchEventsRequestTimeout)
	return &pb.Response{
		Events: evs,
	}, err
}
