package step

// Step is the data for a single endpoint
type Step struct {
	ID        string                 `yaml:"id" validate:"required"`
	Variables map[string]interface{} `yaml:"variables"`
	Request   Request                `yaml:"request" validate:"required"`
	Response  ExpectedResponse       `yaml:"response"`
}

// Request is a single step request data
type Request struct {
	Type                string  `yaml:"type" validate:"required"`
	Endpoint            string  `yaml:"endpoint" validate:"required"`
	Headers             Headers `yaml:"headers"`
	Body                string  `yaml:"body"`
	SkipSSLVerification bool
}

type ExpectedResponse struct {
	Code    *int        `yaml:"code"`
	Headers *Headers    `yaml:"headers"`
	Body    *ExpectBody `yaml:"body"`
}

type Headers map[string]string

type ExpectBody struct {
	Type    *string `yaml:"type"`
	Content *string `yaml:"content"`
	Exact   *bool   `yaml:"exact"`
}
