package step

import (
	"errors"

	"github.com/getapid/apid/log"
	"github.com/getapid/apid/spec"
	"github.com/getapid/apid/template"
	"github.com/getapid/apid/variables"
)

var (
	ErrInterpolationError = errors.New("interpolation_error")
)

type Interpolator struct{}

// NewInterpolator instantiates a new template interpolator
func NewInterpolator() *Interpolator {
	return &Interpolator{}
}

func (i *Interpolator) interpolate(step spec.Step, vars variables.Variables) (spec.Step, error) {
	var err error
	result := step

	// Request interpolation
	if result.Request.URL, err = template.Render(result.Request.URL, vars); err != nil {
		log.L.Errorf("interpolating step URL: %v", err)
		return result, ErrInterpolationError
	}
	if body, err := template.Render(string(result.Request.Body), vars); err != nil {
		log.L.Errorf("interpolating step body: %v", err)
		return result, ErrInterpolationError
	} else {
		result.Request.Body = spec.Body(body)
	}
	headers := make(map[string]string, len(result.Request.Headers))
	for header, value := range result.Request.Headers {
		header, err = template.Render(header, vars)
		if err != nil {
			log.L.Errorf("interpolating step header name: %v", err)
			return result, ErrInterpolationError
		}
		value, err = template.Render(value, vars)
		if err != nil {
			log.L.Errorf("interpolating step header value: %v", err)
			return result, ErrInterpolationError
		}

		headers[header] = value
	}
	result.Request.Headers = headers

	// Response interpolation
	// if result.Expect.Text != nil {
	// 	text, err := template.Render(string(*result.Expect.Text), vars)
	// 	if err != nil {
	// 		log.L.Errorf("interpolating step expected text: %v", err)
	// 		return result, ErrInterpolationError
	// 	}
	// 	result.Expect.Text.Set(text)
	// }

	// if result.Expect.JSON != nil {
	// 	for idx := range result.Expect.JSON {
	// 		if result.Expect.JSON[idx].Is, err = template.Render(result.Expect.JSON[idx].Is, vars); err != nil {
	// 			log.L.Errorf("interpolating step expected json: %v", err)
	// 			return result, ErrInterpolationError
	// 		}
	// 	}
	// }

	// if result.Expect.Headers != nil {
	// 	headerMatchers := make(map[string]string, len(*result.Expect.Headers))
	// 	for header, value := range *result.Expect.Headers {
	// 		header, err := template.Render(string(header), vars)
	// 		if err != nil {
	// 			log.L.Errorf("interpolating step header name: %v", err)
	// 			return result, ErrInterpolationError
	// 		}

	// 		value, err := template.Render(string(value), vars)
	// 		if err != nil {
	// 			log.L.Errorf("interpolating step header value: %v", err)
	// 			return result, ErrInterpolationError
	// 		}

	// 		headerMatchers[header] = value
	// 	}
	// 	result.Expect.Headers.Set(headerMatchers)
	// }

	return result, nil
}
