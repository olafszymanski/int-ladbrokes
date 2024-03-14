package client

import (
	"context"
	"fmt"
	"time"

	"github.com/olafszymanski/int-ladbrokes/internal/config"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/storage"
)

const (
	preMatchEventsUrl = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/EventToOutcomeForClass/%s?simpleFilter=event.isStarted:isFalse&simpleFilter=event.startTime:greaterThanOrEqual:%s&translationLang=en&responseFormat=json&prune=event&prune=market&childCount=event"
	liveEventsUrl     = "https://ss-aka-ori.ladbrokes.com/openbet-ssviewer/Drilldown/2.81/EventToOutcomeForClass/%s?simpleFilter=event.isStarted:isTrue&simpleFilter=event.startTime:greaterThanOrEqual:%s&translationLang=en&responseFormat=json&prune=event&prune=market&childCount=event"
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

	evs, err := c.fetchEvents(
		fmt.Sprintf(liveEventsUrl, cls, time.Now().UTC().Format(time.RFC3339)),
		c.config.LiveEventsRequestTimeout,
	)
	return &pb.Response{
		Events: evs,
	}, err
}

func (c *client) GetPreMatch(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	cls, err := c.getClasses(request.SportType)
	if err != nil {
		return nil, err
	}

	evs, err := c.fetchEvents(
		fmt.Sprintf(preMatchEventsUrl, cls, time.Now().UTC().Format(time.RFC3339)),
		c.config.PreMatchEventsRequestTimeout,
	)
	return &pb.Response{
		Events: evs,
	}, err
}
