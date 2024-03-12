package client

import (
	"bytes"
	"fmt"

	"github.com/olafszymanski/int-ladbrokes/internal/transformer"
	"github.com/olafszymanski/int-sdk/httptls"
	"github.com/olafszymanski/int-sdk/integration/pb"
)

func (c *client) fetchEvents(url string) ([]*pb.Event, error) {
	res, err := c.httpClient.Get(url, httptls.WithTimeout(2000))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequest, err)
	}
	if res.Status != 200 {
		return nil, fmt.Errorf("%w: %v", ErrUnexpectedStatusCode, res.Status)
	}

	r := bytes.NewReader(res.Body)
	return transformer.TransformEvents(r)
}
