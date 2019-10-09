package result

import (
	"github.com/iv-p/apid/common/transaction"
)

type Writer interface {
	Write(transaction.Result)
	Close()
}

type multiWriter struct {
	writers []Writer
}

func NewMultiWriter(w ...Writer) Writer {
	return multiWriter{writers: w}
}

func (w multiWriter) Write(result transaction.Result) {
	for _, writer := range w.writers {
		writer.Write(result)
	}
}

func (w multiWriter) Close() {
	for _, writer := range w.writers {
		writer.Close()
	}
}
