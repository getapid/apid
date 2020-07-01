package step

import (
	"fmt"

	"github.com/getapid/apid-cli/common/template"
	"github.com/getapid/apid-cli/common/variables"
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
		err = fmt.Errorf("interpolating step endpoint: %v", err)
		return PreparedStep(step), err
	}
	if step.Request.Body, err = template.Render(step.Request.Body, vars); err != nil {
		err = fmt.Errorf("interpolating step body: %v", err)
		return PreparedStep(step), err
	}
	if step.Response.Body != nil {
		for idx := range step.Response.Body {
			if step.Response.Body[idx].Is, err = template.Render(step.Response.Body[idx].Is, vars); err != nil {
				err = fmt.Errorf("interpolating step response body: %v", err)
				return PreparedStep(step), err
			}
		}
	}

	headers := make(map[string][]string)
	var key, value string
	for k, vals := range step.Request.Headers {
		if key, err = template.Render(k, vars); err != nil {
			err = fmt.Errorf("interpolating step header key: %v", err)
			return PreparedStep(step), err
		}
		for _, v := range vals {
			if value, err = template.Render(v, vars); err != nil {
				err = fmt.Errorf("interpolating step header value: %v", err)
				return PreparedStep(step), err
			}
			headers[key] = append(headers[key], value)
		}
	}
	step.Request.Headers = headers
	return PreparedStep(step), nil
}
