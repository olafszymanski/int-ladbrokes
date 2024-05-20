package client

import (
	"context"
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-sdk/http"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
	"google.golang.org/protobuf/proto"
)

type client struct {
	config     *config.Config
	httpClient http.Doer
	storage    storage.Storager
	pb.UnimplementedIntegrationServer
}

func NewClient(cfg *config.Config, httpClient http.Doer, storage storage.Storager) pb.IntegrationServer {
	return &client{
		httpClient: httpClient,
		config:     cfg,
		storage:    storage,
	}
}

func (c *client) GetLive(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	evs, err := c.storage.GetHashFields(ctx, fmt.Sprintf(config.LiveEventsStorageKey, request.SportType.String()))
	if err != nil {
		return nil, err
	}
	rsp := &pb.Response{
		Events: make([]*pb.Event, 0, len(evs)),
	}
	ev := &pb.Event{}
	for _, e := range evs {
		proto.Unmarshal(e, ev)
		rsp.Events = append(rsp.Events, ev)
	}
	return rsp, nil
}

func (c *client) GetPreMatch(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	evs, err := c.storage.GetHashFields(ctx, fmt.Sprintf(config.PreMatchEventsStorageKey, request.SportType.String()))
	if err != nil {
		return nil, err
	}
	rsp := &pb.Response{
		Events: make([]*pb.Event, 0, len(evs)),
	}
	ev := &pb.Event{}
	for _, e := range evs {
		proto.Unmarshal(e, ev)
		rsp.Events = append(rsp.Events, ev)
	}
	return rsp, nil
}
