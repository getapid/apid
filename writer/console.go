package writer

import (
	"fmt"
	"io"
	"text/template"

	"github.com/getapid/apid/log"
	"github.com/getapid/apid/step"
)

const templateExecuteMessage = "\n\nsomething went wrong displaying results; try again or report a bug at faq.getapid.com\n\n"

var (
	tmpl *template.Template
)

func init() {
	tmpl = template.New("output").Funcs(template.FuncMap{
		"increment": increment,
		"red":       red,
		"green":     green,
		"indent":    indent,
		"time":      renderTime,
		"add":       add,
		"bold":      bold,
	})

	_, err := tmpl.Parse(schema)
	if err != nil {
		fmt.Println(err)
		tmpl, _ = tmpl.Parse("couldn't compile output template\n")
	}
}

type Result struct {
	Name  string
	Steps []step.Result
}

type Writer interface {
	Prelude()
	Write(Result)
	Conclusion()
}

type Console struct {
	successes, failures int
	silent              bool
	out                 io.Writer
}

func NewConsole(dest io.Writer, silent bool) Writer {
	return &Console{
		silent: silent,
		out:    dest,
	}
}

func (w *Console) Prelude() {
	if w.silent {
		return
	}
	err := tmpl.ExecuteTemplate(w.out, "prelude", nil)
	if err != nil {
		log.L.Debugf("parsing output template: %s", err)
		w.out.Write([]byte(templateExecuteMessage))
	}
}

func (w *Console) Write(r Result) {
	w.count(r)

	log.L.Debugf("writing result to console %v", r)

	if w.silent {
		return
	}

	data := struct {
		Result
	}{r}

	err := tmpl.Execute(w.out, data)
	if err != nil {
		log.L.Debugf("parsing output template: %s", err)
		w.out.Write([]byte(templateExecuteMessage))
	}
}

func (w *Console) count(result Result) {
	if isFailed(result) {
		w.failures++
	} else {
		w.successes++
	}
}

func isFailed(r Result) bool {
	for _, s := range r.Steps {
		if !s.Pass {
			return true
		}
	}
	return false
}

func (w Console) Conclusion() {
	data := struct {
		SuccessSteps, FailedSteps string
	}{fmt.Sprintf("%d", w.successes), fmt.Sprintf("%d", w.failures)}

	err := tmpl.ExecuteTemplate(w.out, "conclusion", data)
	if err != nil {
		log.L.Debugf("parsing output template: %s", err)
		w.out.Write([]byte(templateExecuteMessage))
	}
}
