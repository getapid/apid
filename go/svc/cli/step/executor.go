package step

import (
	"log"
	"time"

	"github.com/iv-p/apid/svc/cli/http"
)

type Executor interface {
	do(Request) HTTPResponse
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

func (e *RequestExecutor) do(request Request) HTTPResponse {
	req := http.Request{
		Method:  request.Type,
		Url:     request.Endpoint,
		Headers: http.Headers(request.Headers),
		Body:    request.Body,
	}
	res, err := e.client.Do(req)
	if err != nil {
		log.Print(err)
	}

	return HTTPResponse{
		StatusCode: res.StatusCode,
		Body:       res.Body,
		Headers:    res.Headers,
		Timings:    transformTimings(res.Timings),
	}
}

func transformTimings(t http.Timings) Timings {
	return Timings{
		DNSLookup:        int64(t.DNSLookup / time.Millisecond),
		TCPConnection:    int64(t.TCPConnection / time.Millisecond),
		TLSHandshake:     int64(t.TLSHandshake / time.Millisecond),
		ServerProcessing: int64(t.ServerProcessing / time.Millisecond),
		ContentTransfer:  int64(t.ContentTransfer / time.Millisecond),
		NameLookup:       int64(t.NameLookup / time.Millisecond),
		Connect:          int64(t.Connect / time.Millisecond),
		PreTransfer:      int64(t.PreTransfer / time.Millisecond),
		StartTransfer:    int64(t.StartTransfer / time.Millisecond),
		Total:            int64(t.Total / time.Millisecond),
	}
}
