package step

//go:generate mockgen -destination=../mock/runner_mock.go -package=mock github.com/getapid/apid/common/step Runner

import (
	"github.com/getapid/apid/common/variables"
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
	extractor    extractor
}

// Result has all the data about the step execution
type Result struct {
	Step     PreparedStep
	Exported Exported
	Valid    ValidationResult
}

func (r Result) OK() bool {
	return r.Valid.OK()
}

// NewRunner instantiates a new HTTPRunner
func NewRunner(
	executor executor,
	validator validator,
	interpolator interpolator,
	extractor extractor) Runner {
	return &runner{executor, validator, interpolator, extractor}
}

// Run interpolates, executes and validates an HTTP step
func (c *runner) Run(step Step, vars variables.Variables) (Result, error) {
	var err error
	var result Result

	result.Step, err = c.interpolator.interpolate(step, vars)
	if err != nil {
		return Result{}, err
	}
	response, err := c.executor.do(result.Step.Request)
	if err != nil {
		return result, err
	}

	result.Valid = c.validator.validate(result.Step.Response, response)
	result.Exported = c.extractor.extract(response, result.Step.Export)
	return result, nil
}
