package step

import (
	"io"
	stdhttp "net/http"
	"testing"

	"github.com/iv-p/apid/common/http"
	"github.com/stretchr/testify/suite"
)

type ValidatorSuite struct {
	suite.Suite
}

func (s *ValidatorSuite) TestValidate() {
	type args struct {
		exp    ExpectedResponse
		actual *http.Response
	}

	incorrectBodyResult := ValidationResult{Errors: map[string]string{"body": ""}}
	incorrectStatusResult := ValidationResult{Errors: map[string]string{"code": ""}}
	incorrectHeadersResult := ValidationResult{Errors: map[string]string{"headers": ""}}
	correctResult := ValidationResult{}

	testCases := [...]struct {
		name      string
		args      args
		expResult ValidationResult
	}{
		{
			name: "status correct",
			args: args{
				exp: ExpectedResponse{
					Code: pint(200),
				},
				actual: &http.Response{
					Response: &stdhttp.Response{
						Body:       &stringReadCloser{},
						StatusCode: 200,
					},
				},
			},
			expResult: correctResult,
		},
		{
			name: "status incorrect",
			args: args{
				exp: ExpectedResponse{
					Code: pint(201),
				},
				actual: &http.Response{
					Response: &stdhttp.Response{
						Body:       &stringReadCloser{},
						StatusCode: 200,
					},
				},
			},
			expResult: incorrectStatusResult,
		},
		{
			name: "header correct",
			args: args{
				exp: ExpectedResponse{
					Headers: &Headers{
						"HEADER1": []string{"value1"},
					},
				},
				actual: &http.Response{
					Response: &stdhttp.Response{
						Body: &stringReadCloser{},
						Header: map[string][]string{
							"HEADER1": {"value1"},
						},
					},
				},
			},
			expResult: correctResult,
		},
		{
			name: "header incorrect",
			args: args{
				exp: ExpectedResponse{
					Headers: &Headers{
						"HEADER1": []string{"value1"},
					},
				},
				actual: &http.Response{
					Response: &stdhttp.Response{
						Body: &stringReadCloser{},
						Header: map[string][]string{
							"HEADER1": {"value2"},
						},
					},
				},
			},
			expResult: incorrectHeadersResult,
		},
		{
			name: "body json exact correct",
			args: args{
				exp: ExpectedResponse{
					Body: &ExpectBody{
						Type:    pstring("json"),
						Content: `{"field1":"exact value"}`,
						Exact:   pbool(true),
					},
				},
				actual: &http.Response{
					Response: &stdhttp.Response{
						Body: &stringReadCloser{},
					},
					ReadBody: []byte(`{"field1":"exact value"}`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body json not exact correct",
			args: args{
				exp: ExpectedResponse{
					Body: &ExpectBody{
						Type:    pstring("json"),
						Content: `{"field1":"value doesn't matter"}`,
						Exact:   pbool(false),
					},
				},
				actual: &http.Response{
					Response: &stdhttp.Response{
						Body: &stringReadCloser{},
					},
					ReadBody: []byte(`{"field1":1}`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body json not exact correct",
			args: args{
				exp: ExpectedResponse{
					Body: &ExpectBody{
						Type:    pstring("plaintext"),
						Content: `hi, what's up`,
						Exact:   pbool(false),
					},
				},
				actual: &http.Response{
					Response: &stdhttp.Response{
						Body: &stringReadCloser{},
					},
					ReadBody: []byte(`hi, what's up mate'`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body json exact incorrect",
			args: args{
				exp: ExpectedResponse{
					Body: &ExpectBody{
						Type:    pstring("json"),
						Content: `{"field1":"value actually matters"}`,
						Exact:   pbool(true),
					},
				},
				actual: &http.Response{
					Response: &stdhttp.Response{
						Body: &stringReadCloser{},
					},
					ReadBody: []byte(`{"field1":"value matters"}`),
				},
			},
			expResult: incorrectBodyResult,
		},
		{
			name: "body json not exact incorrect",
			args: args{
				exp: ExpectedResponse{
					Body: &ExpectBody{
						Type:    pstring("json"),
						Content: `{"field1":{"a":1}}`,
						Exact:   pbool(false),
					},
				},
				actual: &http.Response{
					Response: &stdhttp.Response{
						Body: &stringReadCloser{},
					},
					ReadBody: []byte(`{"field1":{"b":1}}`),
				},
			},
			expResult: incorrectBodyResult,
		},
	}

	validator := httpValidator{}
	for _, t := range testCases {
		actualResult := validator.validate(t.args.exp, t.args.actual)

		s.Truef(keysMatch(actualResult.Errors, t.expResult.Errors), "test case %q", t.name)
		s.Equalf(t.expResult.OK(), actualResult.OK(), "test case %q", t.name)
	}
}

func pint(i int) *int {
	return &i
}

func pstring(s string) *string {
	return &s
}

func pbool(b bool) *bool {
	return &b
}

type stringReadCloser struct {
	body string
}

func (rc *stringReadCloser) Read(dest []byte) (int, error) {
	if len(dest) > len(rc.body) {
		written := copy(dest, rc.body)
		rc.body = ""
		return written, io.EOF
	} else {
		written := copy(dest, rc.body)
		rc.body = rc.body[written:]
		return written, nil
	}
}

func (rc *stringReadCloser) Close() error {
	return nil
}

func keysMatch(this, other map[string]string) bool {
	equal := true
	equal = equal && len(this) == len(other)

	for k := range this {
		_, ok := other[k]
		equal = equal && ok
	}
	return equal
}

func TestValidatorSuite(t *testing.T) {
	suite.Run(t, new(ValidatorSuite))
}
