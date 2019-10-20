package transaction

import (
	"github.com/iv-p/apid/common/result"
	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/variables"
)

type Runner interface {
	Run([]Transaction, variables.Variables) error
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

	Runner
}

func NewTransactionRunner(stepRunner step.Runner, writer result.Writer) Runner {
	return &TransactionRunner{
		stepRunner: stepRunner,
		writer:     writer,
	}
}

func (r *TransactionRunner) Run(transactions []Transaction, vars variables.Variables) error {
	var err error
	for _, transaction := range transactions {
		tVars := variables.New(variables.WithVars(transaction.Variables))
		vars = vars.Merge(tVars)
		res := r.runSingleTransaction(transaction, vars)
		r.writer.Write(res)
	}
	return err
}

func (r *TransactionRunner) runSingleTransaction(transaction Transaction, vars variables.Variables) result.TransactionResult {
	res := result.TransactionResult{
		Steps: make([]step.Result, len(transaction.Steps)),
	}
	for _, step := range transaction.Steps {
		tVars := variables.New(variables.WithVars(step.Variables))
		vars = vars.Merge(tVars)
		result, err := r.stepRunner.Run(step, vars)
		res.Steps = append(res.Steps, result)
		if err != nil {
			break
		}
	}
	return res
}
