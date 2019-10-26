package result

import (
	"encoding/json"
	"fmt"

	"github.com/iv-p/apid/common/transaction"
)

type consoleWriter struct {
	successes, failures int
}

func NewConsoleWriter() Writer {
	return &consoleWriter{}
}

func (w *consoleWriter) Write(result transaction.SingleTransactionResult) {
	if isFailure(result) {
		w.failures++
	} else {
		w.successes++
	}

	bytes, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(fmt.Errorf("couldn't unamrshall result %w", err))
	}
	fmt.Println(string(bytes))
}

func (w *consoleWriter) Close() {
	total := w.failures + w.successes
	fmt.Printf("successful: %d/%d; failed: %d/%d", w.successes, total, w.failures, total)
}

func isFailure(r transaction.SingleTransactionResult) (isFailed bool) {
	for _, s := range r.Steps {
		isFailed = isFailed || !s.OK()
	}
	return
}
