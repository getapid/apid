package variables

import (
	"os"
	"strings"
)

// TODO test this package

type Variables struct {
	data map[string]interface{}
	env  map[string]interface{}
}

// NewFromMap returns a new Variables instance that has the provided map set as
// the main variables namespace and an empty environment namespace
func NewFromMap(m map[string]interface{}) Variables {
	return Variables{
		data: m,
		env:  make(map[string]interface{}),
	}
}

// NewFromEnv takes all the available environment variables and puts them in
// the environment namespace of a new Variables instance
func NewFromEnv() Variables {
	environ := os.Environ()
	v := Variables{
		data: make(map[string]interface{}),
		env:  make(map[string]interface{}, len(environ)),
	}

	for _, e := range environ {
		pair := strings.Split(e, "=")
		v.env[pair[0]] = pair[1]
	}

	return v
}

// Merge another variable instance with this one and return a copy of the result
// not modifying the original set of variables
func (v Variables) Merge(other Variables) Variables {
	return Variables{
		data: mergeMaps(v.data, other.data),
		env:  mergeMaps(v.env, other.env),
	}
}

// Get returns the main namespace of the variables
func (v Variables) Get() map[string]interface{} {
	return v.data
}

// GetEnv returns the environment namespace of the variables
func (v Variables) GetEnv() map[string]interface{} {
	return v.env
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
