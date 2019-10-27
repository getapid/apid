package transaction

import (
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/variables"
)

type Checker interface {
	check(Transaction, variables.Variables) SingleTransactionResult
}

type TransactionChecker struct {
	stepChecker step.Runner

	Checker
}

type SingleTransactionResult struct {
	SequenceIds []string
	Steps       map[string]step.Result
}

func NewStepChecker(stepChecker step.Runner) Checker {
	return &TransactionChecker{
		stepChecker: stepChecker,
	}
}

func (c *TransactionChecker) check(transaction Transaction, vars variables.Variables) SingleTransactionResult {
	res := SingleTransactionResult{
		Steps: make(map[string]step.Result, len(transaction.Steps)),
	}
	for _, step := range transaction.Steps {
		stepVars := variables.New(variables.WithVars(step.Variables))
		vars = vars.Merge(stepVars)
		res.SequenceIds = append(res.SequenceIds, step.ID)
		result := c.stepChecker.Run(step, vars)
		res.Steps[step.ID] = result
		if !result.OK() {
			break
		}
	}
	return res
}
