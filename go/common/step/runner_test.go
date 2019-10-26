package step_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	httpi "github.com/iv-p/apid/common/http"
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/variables"
	"github.com/stretchr/testify/assert"
)

var (
	validResult = step.ValidationResult{Errors: map[string]string{}}
)

func testClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}

func TestHTTPRunner_Check(t *testing.T) {
	vars := variables.New(variables.WithRaw(map[string]interface{}{
		"vars": map[string]interface{}{
			"api-key":  "random-uuid-key",
			"endpoint": "test-endpoint",
		}}),
	)

	type fields struct {
		h http.HandlerFunc
	}
	type args struct {
		step step.Step
		vars variables.Variables
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   step.Result
	}{
		{
			"simple test",
			fields{
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "random-uuid-key", r.Header.Get("X-APID-KEY"))
					assert.Equal(t, "/test-endpoint", r.RequestURI)
					_, _ = w.Write([]byte("OK"))
				}),
			},
			args{
				step.Step{
					Request: step.Request{
						Type:     "GET",
						Endpoint: "http://test.com/{{ vars.endpoint }}",
						Headers: step.Headers{
							"X-APID-KEY": "{{ vars.api-key }}",
						},
					},
				},
				vars,
			},
			step.Result{
				Step: step.PreparedStep{
					Request: step.Request{
						Type:     "GET",
						Endpoint: "http://test.com/test-endpoint",
						Headers: step.Headers{
							"X-APID-KEY": "random-uuid-key",
						},
					},
				},
				Valid: validResult,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, teardown := testClient(tt.fields.h)
			defer teardown()
			timedClient := httpi.NewTimedClient(client)
			c := step.NewRunner(
				step.NewHTTPExecutor(timedClient),
				step.NewHTTPValidator(),
				step.NewTemplateInterpolator())
			got := c.Run(tt.args.step, tt.args.vars)

			assert.Equal(t, tt.want, got, tt.name) // this displays the diffs more nicely
		})
	}
}
