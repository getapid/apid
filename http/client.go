package http

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/getapid/apid/spec"
)

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

	Body    []byte
	Timings Timings
}

// Client adds http request timings as part of the http response
type Client struct {
	client *http.Client
	tracer Tracer
}

// Tracer stores http request timings
type Tracer struct {
	dnsStart,
	dnsDone,
	connectStart,
	connectDone,
	tlsStart,
	tlsDone,
	firstResponseByte,
	wroteRequest,
	finishedRequest time.Time
}

// NewClient creates a default timed client
func NewClient() Client {
	return Client{
		client: http.DefaultClient,
		tracer: Tracer{},
	}
}

// Tracer returns a new httptrace.ClientTrace
func (t *Tracer) Tracer() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { t.dnsStart = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { t.dnsDone = time.Now() },
		ConnectStart:         func(_, _ string) { t.connectStart = time.Now() },
		ConnectDone:          func(net, addr string, err error) { t.connectDone = time.Now() },
		GotFirstResponseByte: func() { t.firstResponseByte = time.Now() },
		TLSHandshakeStart:    func() { t.tlsStart = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { t.tlsDone = time.Now() },
		WroteRequest:         func(_ httptrace.WroteRequestInfo) { t.wroteRequest = time.Now() },
	}
}

// Timings computes and returns the timings for a request
func (t *Tracer) Timings() Timings {
	return Timings{
		DNSLookup:        t.dnsDone.Sub(t.dnsStart),
		TCPConnection:    t.connectDone.Sub(t.connectStart),
		TLSHandshake:     t.tlsDone.Sub(t.tlsStart),
		ServerProcessing: t.firstResponseByte.Sub(t.wroteRequest),
		ContentTransfer:  t.finishedRequest.Sub(t.firstResponseByte),
	}
}

func (t *Tracer) Done() {
	t.finishedRequest = time.Now()
}

// Do executes a http request
func (c Client) Do(ctx context.Context, request spec.Request) (*Response, error) {
	// Prepare request
	r, err := http.NewRequest(request.Type, request.URL, strings.NewReader(request.Body))
	if err != nil {
		log.Fatalf("error preparing request: %s", err)
		return nil, err
	}
	for header, value := range request.Headers {
		r.Header.Add(header, value)
	}
	r = r.WithContext(httptrace.WithClientTrace(ctx, c.tracer.Tracer()))
	r.Close = true
	client := insecure()

	// Execute request
	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	c.tracer.Done()

	// Prepare response
	readBody, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close() // prevents memory leaks
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: response,
		Body:     readBody,
		Timings:  c.tracer.Timings(),
	}, nil
}

func insecure() *http.Client {
	insecureTransport := *http.DefaultTransport.(*http.Transport)
	insecureTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	insecureClient := *http.DefaultClient
	insecureClient.Transport = &insecureTransport
	return &insecureClient
}
