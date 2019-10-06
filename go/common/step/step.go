package step

// Step is the data for a single endpoint
type Step struct {
	ID        string                 `yaml:"id" validate:"string,required"`
	Variables map[string]interface{} `yaml:"variables"`
	Request   Request                `yaml:"request" validate:"struct,required"`
	Response  ExpectedResponse       `yaml:"response"`
}

// Request is a single step request data
type Request struct {
	Type     string  `yaml:"type" validate:"string,required"`
	Endpoint string  `yaml:"endpoint" validate:"string,required"`
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
