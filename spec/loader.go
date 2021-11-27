package spec

import (
	"errors"
	"path/filepath"

	"github.com/getapid/apid/env"
	"github.com/getapid/apid/file"
	"github.com/getapid/apid/log"
)

var (
	ErrNoFilesFound error = errors.New("no_files_found")
	ErrInvalidSpec  error = errors.New("invalid_spec_format")
)

type Loader interface {
	Load(string) []Spec
}

type specLoader struct {
	filereader file.Reader
}

func NewSpecLoader(filereader file.Reader) Loader {
	return &specLoader{filereader: filereader}
}

func (r specLoader) Load(glob string) []Spec {
	files, err := filepath.Glob(glob)
	if err != nil {
		log.L.Errorf("error finding files: %s", err)
	}

	if len(files) == 0 {
		log.L.Fatalf("no files found matching pattern: %s", glob)
	}

	hasError := false

	specs := make(map[string]Spec)
	for _, file := range files {
		jsonSpecs := r.filereader.Load(file, env.LoadVars())
		for name, jsonSpec := range jsonSpecs {
			spec, err := Unmarshal([]byte(jsonSpec))
			hasError = hasError || err != nil

			if _, ok := specs[name]; ok {
				log.L.Fatalf("duplicate spec with name %s", name)
			}

			if len(spec.Name) == 0 {
				spec.Name = name
			}

			specs[name] = spec
		}
	}

	var result []Spec
	for _, spec := range specs {
		result = append(result, spec)
	}

	if hasError {
		log.L.Fatalf("error parsing specs")
	}

	return result
}
