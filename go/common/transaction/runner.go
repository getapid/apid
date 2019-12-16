package transaction

import (
	"fmt"

	"github.com/getapid/apid/common/result"
	"github.com/getapid/apid/common/step"
	"github.com/getapid/apid/common/variables"
)

type Runner interface {
	Run([]Transaction, variables.Variables) bool
}

// Transaction is the definition of a transaction
type Transaction struct {
	ID        string              `yaml:"id" validate:"required"`
	Variables variables.Variables `yaml:"variables"`
	Steps     []step.Step         `yaml:"steps" validate:"required,unique=ID"`
	Matrix    *Matrix             `yaml:"matrix"`
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
	defer r.writer.Close()

	for _, transaction := range transactions {
		tVars := vars.DeepCopy().Merge(transaction.Variables)

		if transaction.Matrix == nil {
			res, ok := r.runSingleTransaction(transaction, tVars)
			r.writer.Write(res)
			allOk = allOk && ok
		} else {
			transactionId := transaction.ID
			run := 1
			for transaction.Matrix.HasNext() {
				transaction.ID = fmt.Sprintf("%s-%d", transactionId, run)
				matrixSet := transaction.Matrix.NextSet()
				matrixTVars := tVars.DeepCopy().Merge(matrixSet)

				res, ok := r.runSingleTransaction(transaction, matrixTVars)
				r.writer.Write(res)
				allOk = allOk && ok
				run++
			}
		}
	}

	return allOk
}

func (r *TransactionRunner) runSingleTransaction(transaction Transaction, vars variables.Variables) (result.TransactionResult, bool) {
	res := result.TransactionResult{Id: transaction.ID}
	exportedVars := variables.New()
	for _, step := range transaction.Steps {
		stepVars := vars.DeepCopy().
			Merge(step.Variables).
			Merge(exportedVars)

		stepResult, _ := r.stepRunner.Run(step, stepVars)
		exportedVars = exportedVars.Merge(variables.New(
			variables.WithRaw(
				map[string]interface{}{
					step.ID: stepResult.Exported.Generic(),
				},
			),
		))
		res.Steps = append(res.Steps, stepResult)
		if !stepResult.OK() {
			return res, false
		}
	}
	return res, true
}
