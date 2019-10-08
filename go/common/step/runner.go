package step

import (
	"github.com/iv-p/apid/common/variables"
)

// Runner is the interface for different types of step runners
type Runner interface {
	Check(Step, variables.Variables) (Result, error)
}

// HTTPRunner interpolates, executes and validates an HTTP step
type HTTPRunner struct {
	executor     Executor
	validator    Validator
	interpolator Interpolator
}

// Result has all the data about the step execution
type Result struct {
	Step  PreparedStep
	Valid ValidationResult
}

// NewHTTPRunner instantiates a new HTTPRunner
func NewHTTPRunner(executor Executor, validator Validator, interpolator Interpolator) Runner {
	return &HTTPRunner{executor, validator, interpolator}
}

// Check interpolates, executes and validates an HTTP step
func (c *HTTPRunner) Check(step Step, vars variables.Variables) (Result, error) {
	prepared, err := c.interpolator.interpolate(step, vars)
	if err != nil {
		return Result{}, err
	}
	response, err := c.executor.do(prepared.Request)
	if err != nil {
		return Result{}, err
	}
	validation, err := c.validator.validate(step.Response, response)
	return Result{prepared, validation}, err
}
