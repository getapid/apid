package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

type TimedClient struct {
	Client
}

type Client interface {
	Do(Request) (Response, error)
}

func NewTimedClient() *TimedClient {
	return &TimedClient{}
}

func (e TimedClient) Do(request Request) (Response, error) {
	var t0, t1, t2, t3, t4, t5, t6, t7 time.Time
	var response Response
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { t1 = time.Now() },
		ConnectStart: func(_, _ string) {
			if t1.IsZero() {
				t1 = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			t2 = time.Now()
		},
		GotConn:              func(_ httptrace.GotConnInfo) { t3 = time.Now() },
		GotFirstResponseByte: func() { t4 = time.Now() },
		TLSHandshakeStart:    func() { t5 = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { t6 = time.Now() },
	}
	bodyReader := bytes.NewReader([]byte(request.Body))
	req, err := http.NewRequest(request.Method, request.Url, bodyReader)
	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	for k, v := range request.Headers {
		req.Header.Add(k, v)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	resp.Body.Close()
	response.Body = string(body)

	response.StatusCode = resp.StatusCode
	t7 = time.Now()
	if t0.IsZero() {
		t0 = t1
	}

	response.Headers = make(Headers)
	for k, vs := range resp.Header {
		for _, v := range vs {
			response.Headers[k] = v
		}
	}

	response.Timings = NewTimings()
	response.Timings.
		WithDNSLookup(t1.Sub(t0)).
		WithTCPConnection(t2.Sub(t1)).
		WithTLSHandshake(t6.Sub(t5)).
		WithServerProcessing(t4.Sub(t3)).
		WithContentTransfer(t7.Sub(t4)).
		WithNameLookup(t1.Sub(t0)).
		WithConnect(t2.Sub(t0)).
		WithPreTransfer(t3.Sub(t0)).
		WithStartTransfer(t4.Sub(t0)).
		WithTotal(t7.Sub(t0))

	return response, nil
}

func isRedirect(resp *http.Response) bool {
	return resp.StatusCode > 299 && resp.StatusCode < 400
}

func readResponseBody(req *http.Request, resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return ""
	}
	return string(body) + ""
}
