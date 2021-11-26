package loader

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/getapid/apid/env"
	"github.com/getapid/apid/log"
	"github.com/google/go-jsonnet"
)

type Loader interface {
	Load(file string, environment env.Vars) (map[string]string, error)
}

type JsonnetLoader struct{}

func (l JsonnetLoader) Load(path string, environment env.Vars) (map[string]string, error) {
	vm := jsonnet.MakeVM()
	vm.StringOutput = true

	for name, value := range environment {
		vm.ExtVar(name, value)
	}

	baseFileName := filepath.Base(path)
	filename := strings.TrimSuffix(baseFileName, filepath.Ext(baseFileName))

	specs, err := vm.EvaluateFileMulti(path)
	if err != nil {
		log.L.Errorf("error loading %s: %s", path, err)
	}

	result := make(map[string]string, len(specs))
	for name, spec := range specs {
		result[fmt.Sprintf("%s::%s", filename, name)] = spec
	}

	return result, err
}
