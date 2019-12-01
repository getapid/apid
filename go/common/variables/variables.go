package variables

import (
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	envNamespace = "env"
	varNamespace = "var"
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

func (v *Variables) UnmarshalYAML(value *yaml.Node) error {
	vars := make(map[string]interface{})
	v.data = map[string]interface{}{
		varNamespace: vars,
	}
	return value.Decode(&vars)
}

func (v Variables) MarshalYAML() (interface{}, error) {
	return v.data[varNamespace], nil
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
		merged := vars.Merge(Variables{data: map[string]interface{}{varNamespace: v}})
		vars.data = merged.data
	}
}

// WithOther places the provided map in the variables namespace of the Variables
func WithOther(v Variables) option {
	return func(vars *Variables) {
		merged := vars.Merge(v)
		vars.data = merged.data
	}
}

// WithRaw places the provided map as the underlying data of the Variables
func WithRaw(v map[string]interface{}) option {
	return func(vars *Variables) {
		vars.data = mergeMaps(vars.data, v)
	}
}

// WithEnv takes all the available environment variables and puts them in
// the environment namespace of the new Variables instance
func WithEnv() option {
	return func(vars *Variables) {
		environ := os.Environ()
		env := make(map[string]interface{}, len(environ)) // we need map[string]interface{} for mergeMaps() to work
		for _, e := range environ {
			pair := strings.SplitN(e, "=", 2)
			env[pair[0]] = pair[1]
		}
		vars.data[envNamespace] = env
	}
}

// Merge another variable instance with this one and return a copy of the result
// not modifying the original set of variables
func (v Variables) Merge(other Variables) Variables {
	return Variables{
		data: mergeMaps(v.data, other.data),
	}
}

// Raw is the underlying map[string]interface{} representation of the vars
func (v Variables) Raw() map[string]interface{} {
	return v.data
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
					continue
				}
			}
		}
		this[key] = deepCopyInterface(newVal)
	}

	return this
}

func (v Variables) DeepCopy() Variables {
	// we are guaranteed this type because the type of v.data is this
	// and deepCopy will return the same type
	copiedMap := deepCopyInterface(v.data).(map[string]interface{})
	return Variables{data: copiedMap}
}

func deepCopyInterface(i interface{}) interface{} {
	return deepCopy(reflect.ValueOf(i)).Interface()
}

func deepCopy(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Interface, reflect.Ptr:
		if value.IsNil() {
			return reflect.ValueOf(nil)
		}
		return deepCopy(value.Elem())
	case reflect.Map:
		copied := reflect.MakeMapWithSize(value.Type(), value.Len())

		iter := value.MapRange()
		for iter.Next() {
			k, v := iter.Key(), iter.Value()
			copied.SetMapIndex(k, deepCopy(v))
		}
		return copied
	case reflect.Slice, reflect.Array:
		copyLen := value.Len()
		copied := reflect.MakeSlice(value.Type(), copyLen, copyLen)

		for i := 0; i < copyLen; i++ {
			copiedElement := deepCopy(value.Index(i))
			copied.Index(i).Set(copiedElement)
		}
		return copied
	default:
		return value
	}
}
