package client

import (
	"context"
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-ladbrokes/internal/storage"
	"github.com/olafszymanski/int-sdk/http"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

type client struct {
	config     *config.Config
	httpClient http.Doer
	storage    *storage.Storage
	pb.UnimplementedIntegrationServer
}

func NewClient(cfg *config.Config, httpClient http.Doer, storage *storage.Storage) pb.IntegrationServer {
	return &client{
		httpClient: httpClient,
		config:     cfg,
		storage:    storage,
	}
}

func (c *client) GetLive(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	evs, err := c.storage.GetEvents(ctx, fmt.Sprintf(config.LiveEventsStorageKey, request.SportType.String()))
	if err != nil {
		return nil, err
	}
	return &pb.Response{
		Events: evs,
	}, nil
}

func (c *client) GetPreMatch(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	evs, err := c.storage.GetEvents(ctx, fmt.Sprintf(config.PreMatchEventsStorageKey, request.SportType.String()))
	if err != nil {
		return nil, err
	}
	return &pb.Response{
		Events: evs,
	}, nil
}
