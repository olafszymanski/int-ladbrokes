package client

import (
	"context"
	"io"

	"github.com/olafszymanski/int-ladbrokes/internal/transformer"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

type TestClientConfig struct {
	Input io.Reader
}

type testClient struct {
	config *TestClientConfig
	pb.UnimplementedIntegrationServer
}

func NewTestClient(config *TestClientConfig) pb.IntegrationServer {
	return &testClient{
		config: config,
	}
}

func (c *testClient) GetLive(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	return nil, nil
}

func (c *testClient) GetPreMatch(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	evs, err := transformer.TransformEvents(c.config.Input)
	return &pb.Response{
		Events: evs,
	}, err
}
