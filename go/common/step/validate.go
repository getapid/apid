package step

import "github.com/iv-p/apid/common/http"

type Validator interface {
	validate(ExpectedResponse, *http.Response) (Validation, error)
}

type HTTPValidator struct {
	Validator
}

type Validation struct {
	OK     bool              // overall check status, true only if every other check passes
	Errors map[string]string // a list of error keys and more information about what caused them
}

func NewHTTPValidator() Validator {
	return &HTTPValidator{}
}

func (v *HTTPValidator) validate(ExpectedResponse, *http.Response) (Validation, error) {
	return Validation{
		OK:     true,
		Errors: make(map[string]string),
	}, nil
}
