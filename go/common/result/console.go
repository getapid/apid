package result

import (
	"encoding/json"
	"fmt"

	"github.com/iv-p/apid/common/transaction"
)

type consoleWriter struct {
}

func NewConsoleWriter() Writer {
	return consoleWriter{}
}

func (w consoleWriter) Write(result transaction.SingleTransactionResult) {
	bytes, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(fmt.Errorf("couldn't unamrshall result %w", err))
	}
	fmt.Println(string(bytes))
}

func (w consoleWriter) Close() {}
