package result

//go:generate mockgen -destination=../mock/console_mock.go -package=mock github.com/iv-p/apid/common/result Writer

import (
	"encoding/json"
	"fmt"

	"github.com/iv-p/apid/common/result"
)

type consoleWriter struct {
	successes, failures int
}

func NewConsoleWriter() result.Writer {
	return &consoleWriter{}
}

func (w *consoleWriter) Write(result result.TransactionResult) {
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

func isFailure(r result.TransactionResult) (isFailed bool) {
	for _, s := range r.Steps {
		isFailed = isFailed || !s.OK()
	}
	return
}
