package writer

import (
	"encoding/json"
	"io"
)

type JSON struct {
	results []Result
	out     io.Writer
}

func NewJSON(dest io.Writer) Writer {
	return &JSON{
		out: dest,
	}
}

func (w *JSON) Prelude() {
}

func (w *JSON) Write(r Result) {
	w.results = append(w.results, r)
}

func (w JSON) Conclusion() {
	pass := 0
	fail := 0
	for _, spec := range w.results {
		ok := true
		for _, step := range spec.Steps {
			ok = ok && step.Pass
		}
		if ok {
			pass++
		} else {
			fail++
		}
	}

	data := struct {
		Passing int      `json:"passing"`
		Failing int      `json:"failing"`
		Steps   []Result `json:"steps"`
	}{
		Passing: pass,
		Failing: fail,
		Steps:   w.results,
	}

	res, _ := json.MarshalIndent(data, "", "  ")
	w.out.Write(res)
}
