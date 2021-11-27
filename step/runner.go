package step

import (
	"github.com/getapid/apid/spec"
	"github.com/getapid/apid/variables"
)

type Runner struct {
	client       HttpClient
	interpolator Interpolator
	matcher      Matcher
	exporter     Exporter
}

type Result struct {
	Name         string
	Pass         bool
	PassedChecks []string
	FailedChecks []string
}

// NewRunner instantiates a new HTTPRunner
func NewRunner(client HttpClient, interpolator Interpolator, matcher Matcher, exporter Exporter) Runner {
	return Runner{client, interpolator, matcher, exporter}
}

// Run interpolates, executes and validates an HTTP step
func (r *Runner) Run(step spec.Step, vars variables.Variables) (Result, variables.Variables, error) {
	var res Result
	s, err := r.interpolator.interpolate(step, vars)
	if err != nil {
		return res, vars, err
	}

	response, err := r.client.Do(s.Request)
	if err != nil {
		return res, vars, err
	}

	ok, passed, failed := r.matcher.validate(step.Expect, response)
	exported := r.exporter.export(response, step.Export)

	res.Name = step.Name
	res.Pass = ok
	res.FailedChecks = failed
	res.PassedChecks = passed
	return res, exported, nil
}
