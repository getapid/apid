package result

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/iv-p/apid/common/result"
	"github.com/iv-p/apid/common/step"
)

var (
	redFail = color.New(color.FgHiRed, color.Bold).Sprint("FAIL")
	greenOk = color.New(color.FgHiGreen, color.Bold).Sprint("OK")
)

type consoleWriter struct {
	successes, failures int
	out                 indentedWriter
}

func NewConsoleWriter(dest io.Writer) result.Writer {
	return &consoleWriter{
		out: indentedWriter{
			out: dest,
		},
	}
}

func (w *consoleWriter) Write(result result.TransactionResult) {
	w.count(result)
	w.out.setIndent(0)

	w.print(result.Id + ":\n")

	w.out.setIndent(4)

	for _, s := range result.Steps {
		if !s.OK() {
			w.printFailedStep(s)
		} else {
			w.printSuccStep(s)
		}
	}
	w.out.setIndent(0)
}

func (w consoleWriter) print(args ...interface{}) {
	_, _ = fmt.Fprint(w.out, args...)
}

func (w consoleWriter) printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w.out, format, args...)
}

func (w *consoleWriter) count(result result.TransactionResult) {
	if isFailed(result) {
		w.failures++
	} else {
		w.successes++
	}
	return
}

func isFailed(r result.TransactionResult) bool {
	for _, s := range r.Steps {
		if !s.OK() {
			return true
		}
	}
	return false
}

func (w consoleWriter) printFailedStep(s step.Result) {
	w.print(redFail + "\t" + s.Step.ID + "\n")

	req := s.Step.Request

	w.printf("request: %s %s\n", req.Type, req.Endpoint)
	w.out.increaseIndent(4)
	if body := formatBody(req); len(body) != 0 {
		w.print(formatBody(req) + "\n")
	}
	w.out.decreaseIndent(4)

	w.print("errors:\n")
	w.out.increaseIndent(4)
	defer w.out.decreaseIndent(4)
	for k, err := range s.Valid.Errors {
		w.print(k + ":\n")
		w.out.increaseIndent(4)
		w.print(err + "\n")
		w.out.decreaseIndent(4)
	}
}

func formatBody(r step.Request) string {
	switch r.Type {
	case "json":
		formatted := &bytes.Buffer{}
		err := json.Indent(formatted, []byte(r.Body), "", "  ")
		if err != nil {
			return r.Body
		}
		return formatted.String()
	default:
		return r.Body
	}
}

func (w consoleWriter) printSuccStep(s step.Result) {
	w.print(greenOk + "\t\t" + s.Step.ID + "\n")
}

func (w consoleWriter) Close() {
	total := w.failures + w.successes
	w.printf("\nsuccessful transactions:\t%d/%d\nfailed transactions:\t\t%d/%d\n", w.successes, total, w.failures, total)
}
