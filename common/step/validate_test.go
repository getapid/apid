package step

import (
	"io"
	"testing"

	"github.com/getapid/apid-cli/common/http"
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
					StatusCode: 200,
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
					StatusCode: 200,
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
					Header: map[string][]string{
						"HEADER1": {"value1"},
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
					Header: map[string][]string{
						"HEADER1": {"value2"},
					},
				},
			},
			expResult: incorrectHeadersResult,
		},
		{
			name: "body json exact correct",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `{"field1":"exact value"}`,
							Subset:   pbool(false),
							KeysOnly: pbool(false),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`{"field1":"exact value"}`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body json exact incorrect",
			args: args{
				exp: ExpectedResponse{

					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `{"field1":"value actually matters"}`,
							Subset:   pbool(false),
							KeysOnly: pbool(false),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`{"field1":"value matters"}`),
				},
			},
			expResult: incorrectBodyResult,
		},
		{
			name: "body not json subset incorrect",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `hi, what's up`,
							Subset:   pbool(false),
							KeysOnly: pbool(false),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`hi, what's up mate'`),
				},
			},
			expResult: incorrectBodyResult,
		},
		{
			name: "body not json subset correct",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `hi, what's up`,
							Subset:   pbool(true),
							KeysOnly: pbool(false),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`hi, what's up mate'`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body json subset correct",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `{"field1":"value doesn't matter"}`,
							Subset:   pbool(false),
							KeysOnly: pbool(true),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`{"field1":1}`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body json subset incorrect",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `{"field1":"value doesn't matter"}`,
							Subset:   pbool(false),
							KeysOnly: pbool(true),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`{"field2":1}`),
				},
			},
			expResult: incorrectBodyResult,
		},
		{
			name: "body nested json subset incorrect",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `{"field1":{"a":1}}`,
							Subset:   pbool(false),
							KeysOnly: pbool(true),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`{"field1":{"b":1}}`),
				},
			},
			expResult: incorrectBodyResult,
		},
		{
			name: "body nested json subset correct",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `{"field1":{"a":1}}`,
							Subset:   pbool(false),
							KeysOnly: pbool(true),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`{"field1":{"a":1}}`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body json array subset incorrect",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `[{"field1":{"a":1}}]`,
							Subset:   pbool(false),
							KeysOnly: pbool(false),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`[{"field1":{"a":1}},{"field1":{"b":1}}]`),
				},
			},
			expResult: incorrectBodyResult,
		},
		{
			name: "body json array subset correct",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `[{"field1":{"a":1}}]`,
							Subset:   pbool(true),
							KeysOnly: pbool(false),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`[{"field1":{"a":1}},{"field1":{"b":1}}]`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body text array subset correct",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `["Tom"]`,
							Subset:   pbool(true),
							KeysOnly: pbool(false),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`["Tom", "Pom", "Nom"]`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body text json subset  correct",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `["Tom"]`,
							Subset:   pbool(true),
							KeysOnly: pbool(false),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`["Tom", "Pom", "Nom"]`),
				},
			},
			expResult: correctResult,
		},
		{
			name: "body json array subset keys only correct",
			args: args{
				exp: ExpectedResponse{
					Body: []*ExpectBody{
						&ExpectBody{
							Is:       `[{"field1":{"a":1}}]`,
							Subset:   pbool(true),
							KeysOnly: pbool(true),
						},
					},
				},
				actual: &http.Response{
					Body: []byte(`[{"field1":{"a":23}},{"field1":{"b":1}}]`),
				},
			},
			expResult: correctResult,
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
