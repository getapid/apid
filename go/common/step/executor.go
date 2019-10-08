package step

import (
	"context"
	"strings"

	"github.com/iv-p/apid/common/http"
)

// Executor is the interface for step executors
type Executor interface {
	do(Request) (*http.Response, error)
}

// HTTPExecutor sends steps as HTTP requests
type HTTPExecutor struct {
	client http.Client
}

// NewHTTPExecutor instantiates a new http executor
func NewHTTPExecutor(client http.Client) Executor {
	return &HTTPExecutor{client: client}
}

func (e *HTTPExecutor) do(request Request) (*http.Response, error) {
	req, err := http.NewRequest(request.Type, request.Endpoint, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}
	return e.client.Do(context.Background(), req)
}
