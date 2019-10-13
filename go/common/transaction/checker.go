package transaction

import (
	"github.com/iv-p/apid/common/http"
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/variables"
)

type Checker interface {
	check(Transaction, variables.Variables) SingleTransactionResult
}

type TransactionChecker struct {
	stepChecker      step.Checker
	stepInterpolator Interpolator

	Checker
}

type StepResult struct {
	Step     step.Step
	Response *http.Response
	Result   step.ValidationResult
}

type SingleTransactionResult struct {
	SequenceIds []string
	Steps       map[string]StepResult
}

func NewStepChecker(stepChecker step.Checker, interpolator Interpolator) Checker {
	return &TransactionChecker{
		stepChecker:      stepChecker,
		stepInterpolator: interpolator,
	}
}

func (c *TransactionChecker) check(transaction Transaction, vars variables.Variables) SingleTransactionResult {
	res := SingleTransactionResult{
		Steps: make(map[string]StepResult),
	}
	for _, step := range transaction.Steps {
		stepVars := variables.NewFromMap(step.Variables)
		vars = vars.Merge(stepVars)
		prepared := c.stepInterpolator.interpolate(step, vars)
		response, result := c.stepChecker.Check(prepared)
		res.SequenceIds = append(res.SequenceIds, prepared.ID)
		res.Steps[prepared.ID] = StepResult{prepared, response, result}
	}
	return res
}
