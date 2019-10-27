package result

import (
	"github.com/iv-p/apid/common/step"
)

type TransactionResult struct {
	Steps []step.Result
}

type Writer interface {
	Write(TransactionResult)
	Close()
}

type multiWriter struct {
	writers []Writer
}

func NewMultiWriter(w ...Writer) Writer {
	return multiWriter{writers: w}
}

func (w multiWriter) Write(result TransactionResult) {
	for _, writer := range w.writers {
		writer.Write(result)
	}
}

func (w multiWriter) Close() {
	for _, writer := range w.writers {
		writer.Close()
	}
}
