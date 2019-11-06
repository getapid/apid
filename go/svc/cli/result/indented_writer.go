package result

import (
	"io"
	"strings"
)

type indentedWriter struct {
	currentIndent string
	out           io.Writer
}

func (w indentedWriter) Write(b []byte) (int, error) {
	str := string(b)
	numNewLines := strings.Count(str, "\n")
	indented := strings.Replace(str, "\n", "\n"+w.currentIndent, numNewLines-1)
	return w.out.Write([]byte(w.currentIndent + indented))
}

func (w *indentedWriter) setIndent(i int) {
	w.currentIndent = strings.Repeat(" ", i)
}

func (w *indentedWriter) increaseIndent(i int) {
	w.setIndent(len(w.currentIndent) + i)
}

func (w *indentedWriter) decreaseIndent(i int) {
	w.setIndent(len(w.currentIndent) - i)
}
