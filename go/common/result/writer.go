package result

import (
	"github.com/iv-p/apid/common/transaction"
)

type Writer interface {
	Write(transaction.Result)
	Close()
}
