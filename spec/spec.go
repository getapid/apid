package spec

import (
	"github.com/getapid/apid/spec/validator"
)

type Spec struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

type Step struct {
	Name    string  `json:"name"`
	Request Request `json:"request"`
	Expect  Expect  `json:"expect"`
	Export  Export  `json:"export"`
}

type Export map[string]string

type Request struct {
	Type    string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type Expect struct {
	Code    *validator.StatusCodeValidator `json:"code"`
	Headers *validator.HeaderValidator     `json:"headers"`
	JSON    []validator.JSONValidator      `json:"json"`
	Text    *validator.BodyTextValidator   `json:"text"`
}
