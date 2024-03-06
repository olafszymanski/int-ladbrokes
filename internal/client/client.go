package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
	"github.com/olafszymanski/int-ladbrokes/internal/mapping"
	"github.com/olafszymanski/int-ladbrokes/internal/transformer"
	"github.com/olafszymanski/int-sdk/integration/pb"
	"github.com/olafszymanski/int-sdk/request"
)

var (
	ErrRequest              = fmt.Errorf("request failed")
	ErrOpenFile             = fmt.Errorf("failed to open file")
	ErrWriteFile            = fmt.Errorf("failed to write to file")
	ErrUnexpectedStatusCode = fmt.Errorf("unexpected status code")
)

type client struct {
	httpClient cycletls.CycleTLS
	pb.UnimplementedIntegrationServer
}

func NewClient() pb.IntegrationServer {
	// TODO: Implement .Start() method
	return &client{
		httpClient: cycletls.Init(),
	}
}

func (c *client) GetLive(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	return nil, nil
}

func (c *client) GetPreMatch(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	cls, err := c.getClasses(mapping.SportTypesCodes[request.SportType])
	if err != nil {
		return nil, err
	}

	clsStr := strings.Join(cls, ",")
	evs, err := c.getEvents(clsStr)
	return &pb.Response{
		Events: evs,
	}, err
}

func (c *client) getClasses(categoryCode int) ([]string, error) {
	url := fmt.Sprintf(classesUrl, categoryCode)
	res, err := request.Do(
		&c.httpClient,
		url,
		http.MethodGet,
		2000,
	)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader([]byte(res.Body))
	return transformer.TransformClasses(r)
}

func (c *client) getEvents(classesIDs string) ([]*pb.Event, error) {
	url := fmt.Sprintf(eventsUrl, classesIDs, time.Now().UTC().Format(time.RFC3339))
	res, err := request.Do(
		&c.httpClient,
		url,
		http.MethodGet,
		2000,
	)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader([]byte(res.Body))
	return transformer.TransformEvents(r)
}
