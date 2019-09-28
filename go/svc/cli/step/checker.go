package step

import "github.com/iv-p/apid/pkg/step"

// - name: login-endpoint-test
// request:
//   type: post
//   endpoint: "https://my.awesome.api/login"
//   headers:
// 	   - X-API-KEY: 1402ed3a-43b7-4b59-8845-be60ede2accc
//   body: '{"username":"john.wick","password": "ultra-secure-password" }'
// expect:
//   code: 200
//   body:
// 	   contains:
// 	     json:
// 		   - key: "jwt"
// 		     value: "\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b"

type Checker interface {
	Check(step.Step) (HTTPResponse, ValidationResult)
}

type StepChecker struct {
	executor  Executor
	validator Validator
}

type Result struct{}

func NewStepChecker(executor Executor, validator Validator) Checker {
	return &StepChecker{executor, validator}
}

func (c *StepChecker) Check(step step.Step) (HTTPResponse, ValidationResult) {
	response := c.executor.do(step.Request)
	result := c.validator.validate(step.Response, response)
	return response, result
}
