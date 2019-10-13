package transaction

import (
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/template"
	"github.com/iv-p/apid/common/variables"
)

type Interpolator interface {
	interpolate(step.Step, variables.Variables) step.Step
}

type StepInterpolator struct {
	Interpolator
}

func NewStepInterpolator() Interpolator {
	return &StepInterpolator{}
}

func (i *StepInterpolator) interpolate(step step.Step, data variables.Variables) step.Step {
	step.Request.Endpoint, _ = template.Render(step.Request.Endpoint, data)
	step.Request.Body, _ = template.Render(step.Request.Body, data)

	headers := make(map[string]string)
	for k, v := range step.Request.Headers {
		key, _ := template.Render(k, data)
		value, _ := template.Render(v, data)
		headers[key] = value
	}
	step.Request.Headers = headers
	return step
}
