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
	if step.Request.Endpoint, err = template.Render(step.Request.Endpoint, vars); err != nil {
		log.L.Warnf("interpolating step endpoint: %v", err)
	}
	if step.Request.Body, err = template.Render(step.Request.Body, vars); err != nil {
		log.L.Warnf("interpolating step body: %v", err)
	}
	if step.Response.Body != nil {
		if step.Response.Body.Content, err = template.Render(step.Response.Body.Content, vars); err != nil {
			log.L.Warnf("interpolating step response body: %v", err)
		}
	}

	headers := make(map[string][]string)
	var key, value string
	for k, vals := range step.Request.Headers {
		if key, err = template.Render(k, vars); err != nil {
			log.L.Warnf("interpolating step header key: %v", err)
		}
		for _, v := range vals {
			if value, err = template.Render(v, vars); err != nil {
				log.L.Warnf("interpolating step header value: %v", err)
			}
			headers[key] = append(headers[key], value)
		}
	}
	step.Request.Headers = headers
	return PreparedStep(step), nil
}
