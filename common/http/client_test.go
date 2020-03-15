package http

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"testing"
	"time"
)

func testingClient(code int, body string, sleep time.Duration) (*http.Client, func()) {
	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(sleep)
		w.WriteHeader(code)
		w.Write([]byte(body))
	})

	s := httptest.NewServer(h)
	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	if sleep > 0*time.Millisecond {
		cli.Timeout = sleep - 5*time.Millisecond
	}
	return cli, s.Close
}

const d = time.Millisecond

type DummyTracer struct{}

func (t *DummyTracer) Tracer() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{}
}

func (t *DummyTracer) Timings() Timings {
	return Timings{
		DNSLookup:        d,
		TCPConnection:    d,
		TLSHandshake:     d,
		ServerProcessing: d,
		ContentTransfer:  d,
	}
}

func (t *DummyTracer) Done() {}

func TestTimedClient_Do(t *testing.T) {
	type fields struct {
		code   int
		body   string
		tracer Tracer
		sleep  time.Duration
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		code    int
		body    []byte
		timings time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    want
		wantErr bool
	}{
		{
			"correct",
			fields{200, "test", &DummyTracer{}, 0 * time.Millisecond},
			args{context.Background()},
			want{200, []byte("test"), d},
			false,
		},
		{
			"timeout",
			fields{200, "test", &DummyTracer{}, 20 * time.Millisecond},
			args{context.Background()},
			want{200, []byte("test"), d},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc, cl := testingClient(tt.fields.code, tt.fields.body, tt.fields.sleep)
			c := TimedClient{
				client:         hc,
				insecureClient: hc,
				tracer:         tt.fields.tracer,
			}

			req, _ := NewRequest("GET", "http://google.com", nil)
			got, err := c.Do(tt.args.ctx, req)

			if (err != nil) != tt.wantErr {
				t.Errorf("TimedClient.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.StatusCode != tt.want.code {
				t.Errorf("TimedClient.Do() = %v, want %v", got.StatusCode, tt.want.code)
			}
			if !bytes.Equal(got.ReadBody, tt.want.body) {
				t.Errorf("TimedClient.Do() = %v, want %v", got.ReadBody, tt.want.body)
			}
			if got.Timings.ContentTransfer != d {
				t.Errorf("TimedClient.Do() = %v, want %v", got.Timings.ContentTransfer, d)
			}
			if got.Timings.DNSLookup != d {
				t.Errorf("TimedClient.Do() = %v, want %v", got.Timings.DNSLookup, d)
			}
			if got.Timings.ServerProcessing != d {
				t.Errorf("TimedClient.Do() = %v, want %v", got.Timings.ServerProcessing, d)
			}
			if got.Timings.TCPConnection != d {
				t.Errorf("TimedClient.Do() = %v, want %v", got.Timings.TCPConnection, d)
			}
			if got.Timings.TLSHandshake != d {
				t.Errorf("TimedClient.Do() = %v, want %v", got.Timings.TLSHandshake, d)
			}
			cl()
		})
	}
}
