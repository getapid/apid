package transaction

import (
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/svc/cli/interpolator"
)

type Interpolator interface {
	interpolate(step.Step, map[string]interface{}) step.Step
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

func (i *StepInterpolator) interpolate(step step.Step, data map[string]interface{}) step.Step {
	step.Request.Endpoint, _ = i.stringInterpolator.Interpolate(step.Request.Endpoint, data)
	step.Request.Body, _ = i.stringInterpolator.Interpolate(step.Request.Body, data)

	headers := make(map[string]string)
	for k, v := range step.Request.Headers {
		key, _ := i.stringInterpolator.Interpolate(k, data)
		value, _ := i.stringInterpolator.Interpolate(v, v)
		headers[key] = value
	}
	step.Request.Headers = headers
	return step
}
