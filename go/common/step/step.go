package step

// Step is the data for a single endpoint
type Step struct {
	ID        string                 `yaml:"id"`
	Variables map[string]interface{} `yaml:"variables"`
	Request   Request                `yaml:"request"`
	Response  ExpectedResponse       `yaml:"response"`
}

// Request is a single step request data
type Request struct {
	Type     string  `yaml:"type"`
	Endpoint string  `yaml:"endpoint"`
	Headers  Headers `yaml:"headers"`
	Body     string  `yaml:"body"`
}

type ExpectedResponse struct {
	Code    *int        `yaml:"code"`
	Headers *Headers    `yaml:"headers"`
	Body    *ExpectBody `yaml:"body"`
}

type Headers map[string]string

type ExpectBody struct {
	Type    *string `yaml:"type"`
	Content *string `yaml:"contains"`
	Exact   *bool   `yaml:"exact"`
}
