package transaction

import (
	"fmt"
	"os"

	"github.com/getapid/apid-cli/common/result"
	"github.com/getapid/apid-cli/common/step"
	"github.com/getapid/apid-cli/common/variables"
)

type Runner interface {
	Run([]Transaction, variables.Variables) bool
}

// Transaction is the definition of a transaction
type Transaction struct {
	ID        string              `yaml:"id" json:"id" validate:"required"`
	Variables variables.Variables `yaml:"variables" json:"variables"`
	Steps     []step.Step         `yaml:"steps" json:"steps" validate:"required,unique=ID"`
	Matrix    *Matrix             `yaml:"matrix" json:"matrix"`
}

type job struct {
	transaction Transaction
	vars        variables.Variables
}

type res struct {
	result  result.TransactionResult
	success bool
}

type TransactionRunner struct {
	stepRunner  step.Runner
	writer      result.Writer
	parallelism int
}

func NewTransactionRunner(stepRunner step.Runner, writer result.Writer, parallelism int) Runner {
	if parallelism < 1 {
		fmt.Printf("parallelism level is negative or zero ( %d ), exiting.\n", parallelism)
		os.Exit(1)
	}
	return &TransactionRunner{
		stepRunner:  stepRunner,
		writer:      writer,
		parallelism: parallelism,
	}
}

func (r *TransactionRunner) Run(transactions []Transaction, vars variables.Variables) bool {
	defer r.writer.Close()

	var jobs []job
	for _, transaction := range transactions {
		tVars := vars.DeepCopy().Merge(transaction.Variables)
		if transaction.Matrix == nil {
			jobs = append(jobs, job{transaction: transaction, vars: tVars})
		} else {
			transactionId := transaction.ID
			run := 0
			for transaction.Matrix.HasNext() {
				transaction.ID = fmt.Sprintf("%s-%d", transactionId, run)
				matrixSet := transaction.Matrix.NextSet()
				matrixTVars := tVars.DeepCopy().Merge(matrixSet)
				jobs = append(jobs, job{transaction: transaction, vars: matrixTVars})
				run++
			}
		}
	}
	jobsChan := make(chan job)
	resultsChan := make(chan res)

	workers := r.parallelism
	if workers > len(jobs) {
		workers = len(jobs)
	}
	for w := 0; w < workers; w++ {
		go r.worker(jobsChan, resultsChan)
	}

	go func() {
		for _, job := range jobs {
			jobsChan <- job
		}
		close(jobsChan)
	}()

	success := true
	for a := 0; a < len(jobs); a++ {
		res := <-resultsChan
		r.writer.Write(res.result)
		success = success && res.success
	}
	close(resultsChan)

	return success
}

func (r *TransactionRunner) worker(jobs <-chan job, results chan<- res) {
	for j := range jobs {
		result, ok := r.runSingleTransaction(j.transaction, j.vars)
		results <- res{result, ok}
	}
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
