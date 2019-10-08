package step

import (
	"github.com/iv-p/apid/common/variables"
)

// Checker is the interface for different types of step checkers
type Checker interface {
	Check(Step, variables.Variables) (Result, error)
}

// HTTPChecker interpolates, executes and validates an HTTP step
type HTTPChecker struct {
	executor     Executor
	validator    Validator
	interpolator Interpolator
}

// Result has all the data about the step execution
type Result struct {
	Step       Step
	Validation Validation
}

// NewHTTPChecker instantiates a new HTTPChecker
func NewHTTPChecker(executor Executor, validator Validator, interpolator Interpolator) Checker {
	return &HTTPChecker{executor, validator, interpolator}
}

// Check interpolates, executes and validates an HTTP step
func (c *HTTPChecker) Check(step Step, vars variables.Variables) (Result, error) {
	prepared, err := c.interpolator.interpolate(step, vars)
	if err != nil {
		return Result{}, err
	}
	response, err := c.executor.do(step.Request)
	if err != nil {
		return Result{}, err
	}
	validation, err := c.validator.validate(step.Response, response)
	return Result{prepared, validation}, err
}
