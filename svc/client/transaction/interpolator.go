package transaction

import (
	"github.com/iv-p/apiping/svc/client/interpolator"
	"github.com/iv-p/apiping/svc/client/step"
	"github.com/iv-p/apiping/svc/client/variables"
)

type Interpolator interface {
	interpolate(step.Step, variables.Variables) step.Step
}

type StepInterpolator struct {
	stringInterpolator interpolator.StringInterpolator
	Interpolator
}

func NewStepInterpolator() Interpolator {
	return &StepInterpolator{
		stringInterpolator: interpolator.NewSimpleStringInterpolator(),
	}
}

func (i *StepInterpolator) interpolate(step step.Step, vars variables.Variables) step.Step {
	vars = vars.Merge(step.Variables)
	var v map[string]interface{} = vars
	step.Request.Endpoint, _ = i.stringInterpolator.Interpolate(step.Request.Endpoint, v)
	step.Request.Body, _ = i.stringInterpolator.Interpolate(step.Request.Body, v)

	var headers map[string]string
	for k, v := range step.Request.Headers {
		key, _ := i.stringInterpolator.Interpolate(k, v)
		value, _ := i.stringInterpolator.Interpolate(v, v)
		headers[key] = value
	}
	step.Request.Headers = headers
	return step
}
