package step

import (
	"github.com/iv-p/apid/common/variables"
)

// Runner takes a step and variables and checks if it
// returns the expected data
type Runner interface {
	Run(Step, variables.Variables) (Result, error)
}

type runner struct {
	executor     executor
	validator    validator
	interpolator interpolator
}

// Result has all the data about the step execution
type Result struct {
	Step  PreparedStep
	Valid ValidationResult
}

// NewRunner instantiates a new HTTPRunner
func NewRunner(executor executor, validator validator, interpolator interpolator) Runner {
	return &runner{executor, validator, interpolator}
}

// Run interpolates, executes and validates an HTTP step
func (c *runner) Run(step Step, vars variables.Variables) (Result, error) {
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
