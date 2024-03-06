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
	"github.com/olafszymanski/int-sdk/integration/server"
	"github.com/olafszymanski/int-sdk/request"
	"github.com/sirupsen/logrus"
)

var (
	ErrRequest              = fmt.Errorf("request failed")
	ErrOpenFile             = fmt.Errorf("failed to open file")
	ErrWriteFile            = fmt.Errorf("failed to write to file")
	ErrUnexpectedStatusCode = fmt.Errorf("unexpected status code")
)

type Ladbrokes struct {
	httpClient cycletls.CycleTLS
	logger     *logrus.Entry
	pb.UnimplementedIntegrationServer
}

func New(logger *logrus.Entry) pb.IntegrationServer {
	server.Start("8080")
	return &Ladbrokes{
		httpClient: cycletls.Init(),
		logger:     logger,
	}
}

func (c *Ladbrokes) GetLive(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	return nil, nil
}

func (c *Ladbrokes) GetPreMatch(ctx context.Context, request *pb.Request) (*pb.Response, error) {
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

func (c *Ladbrokes) getClasses(categoryCode int) ([]string, error) {
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
	return transformer.UnmarshallClasses(r)
}

func (c *Ladbrokes) getEvents(classesIDs string) ([]*pb.Event, error) {
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
	return transformer.UnmarshallEvents(r)
}
