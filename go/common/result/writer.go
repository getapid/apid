package result

//go:generate mockgen -destination=../mock/writer_mock.go -package=mock github.com/getapid/apid/common/result Writer

import (
	"github.com/getapid/apid/common/step"
)

type TransactionResult struct {
	Id    string
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
