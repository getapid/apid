package step

import (
	"context"
	"strings"

	"github.com/getapid/apid-cli/common/http"
	"github.com/getapid/apid-cli/common/log"
)

type executor interface {
	do(Request) (*http.Response, error)
}

type httpExecutor struct {
	client http.Client
}

// NewHTTPExecutor instantiates a new http executor
func NewHTTPExecutor(client http.Client) executor {
	return &httpExecutor{client: client}
}

func (e *httpExecutor) do(request Request) (*http.Response, error) {
	req, err := http.NewRequest(request.Type, request.Endpoint, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}

	req.SkipVerify = request.SkipSSLVerification
	for k, v := range request.Headers {
		for _, subV := range v {
			req.Header.Add(k, subV)
		}
	}

	log.L.Debugw("executing request",
		"method", req.Method,
		"endpoint", req.URL,
		"headers", req.Header,
		"body", request.Body,
	)

	resp, err := e.client.Do(context.Background(), req)
	if err != nil {
		return nil, err
	}

	log.L.Debugw("received response",
		"endpoint", req.URL,
		"method", req.Method,
		"body", string(resp.ReadBody),
		"headers", resp.Header,
	)

	return resp, nil
}
