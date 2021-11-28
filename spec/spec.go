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
	Body    Body              `json:"body"`
}

type Body string

func (b *Body) UnmarshalJSON(data []byte) error {
	d := string(data)
	*b = Body(d)
	return nil
}

type Expect struct {
	Code    *validator.StatusCodeValidator `json:"code"`
	Headers *validator.HeaderValidator     `json:"headers"`
	Body    *validator.BodyValidator       `json:"body"`
}
