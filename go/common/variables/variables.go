package variables

import (
	"os"
	"strings"
)

type Variables struct {
	data map[string]interface{}
}

// New returns a new Variables instance that has the provided map set as
// the main variables namespace and an empty environment namespace
func New(opts ...option) Variables {
	v := newEmptyVars()
	for _, o := range opts {
		o(&v)
	}
	return v
}

func newEmptyVars() Variables {
	return Variables{
		data: make(map[string]interface{}),
	}
}

type option func(variables *Variables)

// WithVars places the provided map in the variables namespace of the Variables
func WithVars(v map[string]interface{}) option {
	return func(vars *Variables) {
		vars.data["variables"] = v
	}
}

// WithRaw places the provided map as the underlying data of the Variables
func WithRaw(v map[string]interface{}) option {
	return func(vars *Variables) {
		vars.data = v
	}
}

// WithEnv takes all the available environment variables and puts them in
// the environment namespace of the new Variables instance
func WithEnv() option {
	return func(vars *Variables) {
		environ := os.Environ()
		env := make(map[string]interface{}, len(environ)) // we need map[string]interface{} for mergeMaps() to work
		for _, e := range environ {
			pair := strings.Split(e, "=")
			env[pair[0]] = pair[1]
		}
		vars.data["env"] = env
	}
}

// Merge another variable instance with this one and return a copy of the result
// not modifying the original set of variables
func (v Variables) Raw() map[string]interface{} {
	return v.data
}

func (v Variables) Merge(other Variables) Variables {
	return Variables{
		data: mergeMaps(v.data, other.data),
	}
}

func mergeMaps(this, other map[string]interface{}) map[string]interface{} {
	if this == nil {
		return other
	}
	for key, newVal := range other {
		if existingVal, ok := this[key]; ok {
			// if the existing value isn't mergable we skip it
			if existingMap, ok := existingVal.(map[string]interface{}); ok {
				// if the new value isn't mergable we skip it
				if newMap, ok := newVal.(map[string]interface{}); ok {
					this[key] = mergeMaps(existingMap, newMap)
				}
			}
		} else {
			this[key] = newVal
		}
	}
	return this
}
