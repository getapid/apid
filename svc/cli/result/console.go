package result

import (
	"fmt"
	"io"
	"text/template"

	"github.com/fatih/color"
	"github.com/getapid/cli/common/log"
	"github.com/getapid/cli/common/result"
)

const templateExecuteMessage = "something went wrong displaying results; try again or report a bug at faq.getapid.com"

var (
	redFail = color.New(color.FgHiRed, color.Bold).Sprint("FAIL")
	greenOk = color.New(color.FgHiGreen, color.Bold).Sprint("OK")
	tmpl    *template.Template
)

func init() {
	tmpl = template.New("output").Funcs(template.FuncMap{
		"increment": increment,
		"greenOk":   func() string { return greenOk },
		"redFail":   func() string { return redFail },
		"indent":    indent,
		"time":      renderTime,
		"add":       add,
	})

	_, err := tmpl.Parse(schema)
	if err != nil {
		fmt.Println(err)
		tmpl, _ = tmpl.Parse("couldn't compile output template\n")
	}
}

type consoleWriter struct {
	successes, failures int
	showTimings         bool
	out                 io.Writer
}

func NewConsoleWriter(dest io.Writer, displayTimings bool) result.Writer {
	return &consoleWriter{
		out:         dest,
		showTimings: displayTimings,
	}
}

func (w *consoleWriter) Write(r result.TransactionResult) {
	w.count(r)

	data := struct {
		result.TransactionResult
		ShowTimings bool
	}{r, w.showTimings}

	err := tmpl.Execute(w.out, data)
	if err != nil {
		log.L.Debugf("parsing output template: %s", err)
		fmt.Println(templateExecuteMessage)
	}
}

func (w *consoleWriter) count(result result.TransactionResult) {
	if isFailed(result) {
		w.failures++
	} else {
		w.successes++
	}
}

func isFailed(r result.TransactionResult) bool {
	for _, s := range r.Steps {
		if !s.OK() {
			return true
		}
	}
	return false
}

func (w consoleWriter) Close() {
	data := struct {
		SuccessSteps, FailedSteps int
	}{w.successes, w.failures}

	err := tmpl.ExecuteTemplate(w.out, "closingLines", data)
	if err != nil {
		log.L.Debugf("parsing output template: %s", err)
		fmt.Println(templateExecuteMessage)
	}
}
