package step

import (
	"github.com/iv-p/apid/common/template"
	"github.com/iv-p/apid/common/variables"
)

// Interpolator is the interface for different types of step interpolators
type Interpolator interface {
	interpolate(Step, variables.Variables) (Step, error)
}

// TemplateInterpolator uses the template package to interpolate a step
type TemplateInterpolator struct{}

// NewTemplateInterpolator instantiates a new template interpolator
func NewTemplateInterpolator() *TemplateInterpolator {
	return &TemplateInterpolator{}
}

func (i *TemplateInterpolator) interpolate(step Step, vars variables.Variables) (Step, error) {
	step.Request.Endpoint, _ = template.Render(step.Request.Endpoint, vars.Get())
	step.Request.Body, _ = template.Render(step.Request.Body, vars.Get())

	headers := make(map[string]string)
	for k, v := range step.Request.Headers {
		key, _ := template.Render(k, vars.Get())
		value, _ := template.Render(v, v)
		headers[key] = value
	}
	step.Request.Headers = headers
	return step, nil
}
