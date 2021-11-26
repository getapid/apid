package step

import (
	"context"

	"github.com/getapid/apid/http"
	"github.com/getapid/apid/log"
	"github.com/getapid/apid/spec"
)

type HttpClient struct {
	httpClient http.Client
}

// NewHTTPClient instantiates a new http executor
func NewHTTPClient(client http.Client) HttpClient {
	return HttpClient{httpClient: client}
}

func (c *HttpClient) Do(request spec.Request) (*http.Response, error) {
	log.L.Debugw("executing request",
		"method", request.Type,
		"endpoint", request.URL,
		"headers", request.Headers,
		"body", request.Body,
	)

	resp, err := c.httpClient.Do(context.Background(), request)
	if err != nil {
		return nil, err
	}

	log.L.Debugw("received response",
		"endpoint", request.URL,
		"method", request.Type,
		"body", string(resp.Body),
		"headers", resp.Header,
	)

	return resp, nil
}
