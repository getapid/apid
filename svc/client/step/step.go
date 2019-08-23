package step

import "github.com/iv-p/apiping/svc/client/variables"

// Step is the data for a single endpoint
type Step struct {
	ID        string              `yaml:"id"`
	Variables variables.Variables `yaml:"variables"`
	Request   Request             `yaml:"request"`
	Response  ExpectedResponse    `yaml:"response"`
}

// Request is a single step request data
type Request struct {
	Type     string            `yaml:"type"`
	Endpoint string            `yaml:"endpoint"`
	Headers  map[string]string `yaml:"headers"`
	Body     string            `yaml:"body"`
}

type ExpectedResponse struct {
	Code    *int             `yaml:"code"`
	Headers *ExpectedHeaders `yaml:"headers"`
	Body    *ExpectBody      `yaml:"body"`
}

type ExpectedHeaders map[string]string

type ExpectBody struct {
	Type    *string `yaml:"type"`
	Content *string `yaml:"contains"`
	Exact   *bool   `yaml:"exact"`
}
