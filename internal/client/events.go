package client

import (
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/transformer"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

func (c *client) fetchEvents(url string, timeout int) ([]*pb.Event, error) {
	res, err := c.httpClient.Get(url, httptls.WithTimeout(timeout))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}
	return transformer.TransformEvents(res.Body)
}
