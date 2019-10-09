package result

import (
	"github.com/iv-p/apid/common/transaction"
)

// Writer is the interface for result writers
type Writer interface {
	Write(transaction.Result)
	Close()
}
