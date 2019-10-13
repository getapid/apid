package transaction

import (
	"github.com/iv-p/apid/common/variables"
)

type Service interface {
	Check([]Transaction, variables.Variables) MultipleTransactionsResults
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

func (s *TransactionService) Check(transactions []Transaction, vars variables.Variables) MultipleTransactionsResults {
	res := make(map[string]SingleTransactionResult)
	for _, transaction := range transactions {
		tVars := variables.New(variables.WithVars(transaction.Variables))
		vars = vars.Merge(tVars)
		res[transaction.ID] = s.checker.check(transaction, vars)
	}
	return res
}
