package transaction

import (
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/variables"
)

type Runner interface {
	run(Transaction, variables.Variables) SingleTransactionResult
}

type TransactionRunner struct {
	stepRunner step.Runner

	Runner
}

type SingleTransactionResult struct {
	Steps []step.Result
}

func NewTransactionRunner(stepRunner step.Runner) Runner {
	return &TransactionRunner{
		stepRunner: stepRunner,
	}
}

func (c *TransactionRunner) run(transaction Transaction, vars variables.Variables) SingleTransactionResult {
	res := SingleTransactionResult{
		Steps: make([]step.Result, len(transaction.Steps)),
	}
	for _, step := range transaction.Steps {
		tVars := variables.New(variables.WithVars(step.Variables))
		vars = vars.Merge(tVars)
		result, err := c.stepRunner.Run(step, vars)
		res.Steps = append(res.Steps, result)
		if err != nil {
			break
		}
	}
	return res
}
