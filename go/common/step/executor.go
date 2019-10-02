package step

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/iv-p/apid/common/http"
)

type Executor interface {
	do(Request) *http.Response
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

func (e *RequestExecutor) do(request Request) *http.Response {
	req, err := http.NewRequest(request.Type, request.Endpoint, strings.NewReader(request.Body))
	for k, v := range request.Headers {
		req.Header.Set(k, v)
	}
	res, err := e.client.Do(context.Background(), req)
	if err != nil {
		log.Print(err)
	}

	return res
}

func transformTimings(t http.Timings) Timings {
	return Timings{
		DNSLookup:        int64(t.DNSLookup / time.Millisecond),
		TCPConnection:    int64(t.TCPConnection / time.Millisecond),
		TLSHandshake:     int64(t.TLSHandshake / time.Millisecond),
		ServerProcessing: int64(t.ServerProcessing / time.Millisecond),
		ContentTransfer:  int64(t.ContentTransfer / time.Millisecond),
	}
}
