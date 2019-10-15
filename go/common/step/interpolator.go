package step

import (
	"github.com/iv-p/apid/common/log"
	"github.com/iv-p/apid/common/template"
	"github.com/iv-p/apid/common/variables"
)

type interpolator interface {
	interpolate(Step, variables.Variables) (PreparedStep, error)
}

// PreparedStep is the same as step, but with replaced template tokens
type PreparedStep Step

type templateInterpolator struct{}

// NewTemplateInterpolator instantiates a new template interpolator
func NewTemplateInterpolator() *templateInterpolator {
	return &templateInterpolator{}
}

func (i *templateInterpolator) interpolate(step Step, vars variables.Variables) (PreparedStep, error) {
	var err error
	if step.Request.Endpoint, err = template.Render(step.Request.Endpoint, vars.Get()); err != nil {
		log.L.Warnf("interpolating step endpoint: %v", err)
	}
	if step.Request.Body, err = template.Render(step.Request.Body, vars.Get()); err != nil {
		log.L.Warnf("interpolating step body: %v", err)
	}

	headers := make(map[string]string)
	var key, value string
	for k, v := range step.Request.Headers {
		if key, err = template.Render(k, vars.Get()); err != nil {
			log.L.Warnf("interpolating step header key: %v", err)
		}
		if value, err = template.Render(v, vars.Get()); err != nil {
			log.L.Warnf("interpolating step header value: %v", err)
		}
		headers[key] = value
	}
	step.Request.Headers = headers
	return PreparedStep(step), nil
}
