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
		Steps: make(map[string]step.Result),
	}
	for _, step := range transaction.Steps {
		vars = vars.Merge("variables", step.Variables)
		res.SequenceIds = append(res.SequenceIds, step.ID)
		result, err := c.stepChecker.Check(step, vars)
		res.Steps[step.ID] = result
		if err != nil {
			break
		}
	}
	return res
}
