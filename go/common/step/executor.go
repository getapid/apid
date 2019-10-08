package step

import (
	"context"
	"strings"

	"github.com/iv-p/apid/common/http"
)

type Executor interface {
	do(Request) (*http.Response, error)
}

type RequestExecutor struct {
	client http.Client

	Executor
}

type HTTPResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
	Timings    Timings
}

type Timings struct {
	DNSLookup,
	TCPConnection,
	TLSHandshake,
	ServerProcessing,
	ContentTransfer,
	NameLookup,
	Connect,
	PreTransfer,
	StartTransfer,
	Total int64
}

func NewRequestExecutor(client http.Client) Executor {
	return &RequestExecutor{client: client}
}

func (e *RequestExecutor) do(request Request) (*http.Response, error) {
	req, err := http.NewRequest(request.Type, request.Endpoint, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}
	return e.client.Do(context.Background(), req)
}
