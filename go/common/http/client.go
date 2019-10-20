package http

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"
)

// Request is a http request
type Request struct {
	*http.Request

	SkipVerify bool
}

// Timings holds the timings data for a http request
type Timings struct {
	DNSLookup,
	TCPConnection,
	TLSHandshake,
	ServerProcessing,
	ContentTransfer time.Duration
}

// Response a http response with added timing information
type Response struct {
	*http.Response
	Timings Timings
}

// Client is the interface of a http client
type Client interface {
	Do(context.Context, *Request) (*Response, error)
}

// TimedClient adds http request timings as part of the http response
type TimedClient struct {
	client         *http.Client
	insecureClient *http.Client
	tracer         Tracer
}

// Tracer is the interface for a tracer
type Tracer interface {
	Tracer() *httptrace.ClientTrace
	Timings() Timings
}

// DefaultClient is the default HTTP client
var DefaultClient = http.DefaultClient

// DefaultTracer stores http request timings
type DefaultTracer struct {
	dnsStart,
	dnsDone,
	connectStart,
	connectDone,
	tlsStart,
	tlsDone,
	firstResponseByte,
	wroteRequest time.Time
}

// NewTimedClient creates a default timed client
func NewTimedClient(client *http.Client) *TimedClient {
	insecureClient := insecure(client)
	return &TimedClient{
		client:         client,
		insecureClient: insecureClient,
		tracer:         &DefaultTracer{},
	}
}

// Tracer returns a new httptrace.ClientTrace
func (t *DefaultTracer) Tracer() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { t.dnsStart = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { t.dnsDone = time.Now() },
		ConnectStart:         func(_, _ string) { t.connectStart = time.Now() },
		ConnectDone:          func(net, addr string, err error) { t.connectDone = time.Now() },
		GotFirstResponseByte: func() { t.firstResponseByte = time.Now() },
		TLSHandshakeStart:    func() { t.tlsStart = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { t.tlsDone = time.Now() },
	}
}

// Timings computes and returns the timings for a request
func (t *DefaultTracer) Timings() Timings {
	return Timings{
		DNSLookup:        t.dnsStart.Sub(t.dnsDone),
		TCPConnection:    t.connectStart.Sub(t.connectDone),
		TLSHandshake:     t.tlsStart.Sub(t.tlsDone),
		ServerProcessing: t.firstResponseByte.Sub(t.connectDone),
		ContentTransfer:  t.wroteRequest.Sub(t.firstResponseByte),
	}
}

// Do executes a http request
func (c TimedClient) Do(ctx context.Context, req *Request) (*Response, error) {
	var res = &Response{}
	var err error
	req.Request = req.WithContext(httptrace.WithClientTrace(ctx, c.tracer.Tracer()))
	// Should we log the error here, or propagate upwards and ingest quietly?
	client := c.client
	if req.SkipVerify {
		client = c.insecureClient
	}
	res.Response, err = client.Do(req.Request)
	res.Timings = c.tracer.Timings()
	return res, err
}

func insecure(source *http.Client) *http.Client {
	insecureTransport := *http.DefaultTransport.(*http.Transport)
	insecureTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	insecureClient := *source
	insecureClient.Transport = &insecureTransport

	return &insecureClient
}

// NewRequest is a helper to quickly create a new http request
func NewRequest(verb string, url string, body io.Reader) (*Request, error) {
	r := &Request{}
	var err error
	r.Request, err = http.NewRequest(verb, url, body)
	return r, err
}
