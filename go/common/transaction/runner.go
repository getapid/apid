package transaction

import (
	"github.com/iv-p/apid/common/result"
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/variables"
)

type Runner interface {
	Run([]Transaction, variables.Variables) bool
}

// Transaction is the definition of a transaction
type Transaction struct {
	ID        string                 `yaml:"id" validate:"required"`
	Variables map[string]interface{} `yaml:"variables"`
	Steps     []step.Step            `yaml:"steps" validate:"required,unique=ID"`
}

type TransactionRunner struct {
	stepRunner step.Runner
	writer     result.Writer
}

func NewTransactionRunner(stepRunner step.Runner, writer result.Writer) Runner {
	return &TransactionRunner{
		stepRunner: stepRunner,
		writer:     writer,
	}
}

func (r *TransactionRunner) Run(transactions []Transaction, vars variables.Variables) bool {
	allOk := true
	for _, transaction := range transactions {
		tVars := variables.New(variables.WithVars(transaction.Variables))
		vars = vars.Merge(tVars)
		res, ok := r.runSingleTransaction(transaction, vars)
		r.writer.Write(res)
		allOk = allOk && ok
	}
	return allOk
}

func (r *TransactionRunner) runSingleTransaction(transaction Transaction, vars variables.Variables) (result.TransactionResult, bool) {
	ok := true
	res := result.TransactionResult{}
	for _, step := range transaction.Steps {
		stepVars := variables.New(variables.WithVars(step.Variables))
		vars = vars.Merge(stepVars)
		stepResult := r.stepRunner.Run(step, vars)
		res.Steps = append(res.Steps, stepResult)
		if !stepResult.OK() {
			return res, false
		}
	}
	return res, ok
}
