package step

import (
	"context"
	"strings"

	"github.com/iv-p/apid/common/http"
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
	return e.client.Do(context.Background(), req)
}
