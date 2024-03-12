package client

import (
	"context"

	"github.com/olafszymanski/int-ladbrokes/internal/transformer"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

type TestClientConfig struct {
	LiveInput     []byte
	PreMatchInput []byte
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
	evs, err := transformer.TransformEvents(c.config.LiveInput)
	return &pb.Response{
		Events: evs,
	}, err
}

func (c *testClient) GetPreMatch(ctx context.Context, request *pb.Request) (*pb.Response, error) {
	evs, err := transformer.TransformEvents(c.config.PreMatchInput)
	return &pb.Response{
		Events: evs,
	}, err
}
