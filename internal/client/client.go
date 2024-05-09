package client

import (
	"context"
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
	"google.golang.org/protobuf/proto"
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
	evs, err := c.storage.GetHashFields(ctx, fmt.Sprintf("LIVE_EVENTS_%s", request.SportType.String()))
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
	evs, err := c.storage.GetHashFields(ctx, fmt.Sprintf("PRE_MATCH_EVENTS_%s", request.SportType.String()))
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
