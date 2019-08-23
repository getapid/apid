package transaction

import "github.com/iv-p/apiping/svc/client/variables"

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

func (s *TransactionService) Check(transactions []Transaction, variables variables.Variables) MultipleTransactionsResults {
	res := make(map[string]SingleTransactionResult)
	for _, transaction := range transactions {
		vars := variables.Merge(transaction.Variables)
		res[transaction.ID] = s.checker.check(transaction, vars)
	}
	return res
}
