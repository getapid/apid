package transaction

import (
	"github.com/iv-p/apid/common/transaction"
	"github.com/iv-p/apid/svc/cli/variables"
)

type Service interface {
	Check([]transaction.Transaction, variables.Variables) MultipleTransactionsResults
}

type TransactionService struct {
	checker Checker

	Service
}

type MultipleTransactionsResults map[string]SingleTransactionResult

func NewTransactionService(checker Checker) Service {
	return &TransactionService{
		checker: checker,
	}
}

func (s *TransactionService) Check(transactions []transaction.Transaction, vars variables.Variables) MultipleTransactionsResults {
	res := make(map[string]SingleTransactionResult)
	for _, transaction := range transactions {
		vars = vars.Merge("variables", transaction.Variables)
		res[transaction.ID] = s.checker.check(transaction, vars)
	}
	return res
}
