package step_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	httpi "github.com/iv-p/apid/common/http"
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/variables"
	"github.com/stretchr/testify/assert"
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
	vars := variables.Variables{
		Data: map[string]interface{}{
			"vars": map[string]interface{}{
				"api-key":  "random-uuid-key",
				"endpoint": "test-endpoint",
			},
		},
	}

	type fields struct {
		h http.HandlerFunc
	}
	type args struct {
		step step.Step
		vars variables.Variables
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    step.Result
		wantErr bool
	}{
		{
			"simple test",
			fields{
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "random-uuid-key", r.Header.Get("X-APID-KEY"))
					assert.Equal(t, "/test-endpoint", r.RequestURI)
					w.Write([]byte("OK"))
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
				step.PreparedStep{
					Request: step.Request{
						Type:     "GET",
						Endpoint: "http://test.com/test-endpoint",
						Headers: step.Headers{
							"X-APID-KEY": "random-uuid-key",
						},
					},
				},
				step.ValidationResult{
					true,
					map[string]string{},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, teardown := testClient(tt.fields.h)
			defer teardown()
			timedClient := httpi.NewTimedClient(client)
			c := step.NewHTTPRunner(
				step.NewHTTPExecutor(timedClient),
				step.NewHTTPValidator(),
				step.NewTemplateInterpolator())
			got, err := c.Check(tt.args.step, tt.args.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTTPRunner.Check() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HTTPRunner.Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
