package config

import (
	"fmt"
	"reflect"

	"github.com/iv-p/apid/common/step"
	"go.uber.org/multierr"
)

// DefaultValidator does not perform any validations.
type DefaultValidator struct {
}

// Validate method performs validation and returns result and optional error.
func (v DefaultValidator) Validate(val interface{}) (bool, error) {
	return true, nil
}

// VersionValidator validates string if it is a valid version
type VersionValidator struct{}

// Validate method performs validation and returns result and optional error.
func (v VersionValidator) Validate(val interface{}) (bool, error) {
	s, ok := val.(string)

	if !ok {
		return false, fmt.Errorf("version must be a string")
	}

	if s != "1" {
		return false, fmt.Errorf("supported versions: \"1\"")
	}

	return true, nil
}

// UniqueValidator validates all elements in the slice have unique property (defined by the prop)
type UniqueValidator struct {
	prop string
}

// Validate method performs validation and returns result and optional error.
func (v UniqueValidator) Validate(val interface{}) (b bool, err error) {
	if val == nil {
		return true, nil
	}

	slice := reflect.ValueOf(val)

	if slice.Kind() != reflect.Slice {
		return false, fmt.Errorf("must be a slice")
	}

	seen := make(map[string]bool)
	for i := 0; i < slice.Len(); i++ {
		valueInQuestion := slice.Index(i).FieldByName(v.prop).String()
		if _, ok := seen[valueInQuestion]; ok {
			err = multierr.Append(err, fmt.Errorf("should contain unique values for %s. Found multiple %s", v.prop, valueInQuestion))
		}
		seen[valueInQuestion] = true
	}

	return err == nil, err
}

// RequiredValidator validates slices and their length
type RequiredValidator struct{}

// Validate method performs validation and returns result and optional error.
func (v RequiredValidator) Validate(val interface{}) (b bool, err error) {
	if val == nil {
		return false, fmt.Errorf("must not be nil")
	}

	thing := reflect.ValueOf(val)

	if thing.Kind() != reflect.Struct && thing.Len() == 0 {
		return false, fmt.Errorf("length must not be 0")
	}

	return true, nil
}

// ExpectBodyValidator validates a step.ExpectBody so that the type and the exact fields make sense together.
// It also sets the defaults for type and exact
type ExpectBodyValidator struct{}

func (v ExpectBodyValidator) Validate(val interface{}) (b bool, err error) {
	if val == nil || (reflect.ValueOf(val).Kind() == reflect.Ptr && reflect.ValueOf(val).IsNil()) {
		return true, nil
	}
	expBody, ok := val.(*step.ExpectBody)
	if !ok {
		return false, fmt.Errorf("must be *step.ExpectBody")
	}

	const (
		typeJson  = "json"
		typePlain = "plaintext"
	)

	typ := typePlain
	if expBody.Type != nil {
		typ = *expBody.Type
	} else {
		(*expBody).Type = &typ
	}

	exact := true
	if expBody.Exact != nil {
		exact = *expBody.Exact
	} else {
		(*expBody).Exact = &exact
	}

	switch typ {
	case typeJson, typePlain:
	default:
		return false, fmt.Errorf("unsupported content type %q", typ)
	}

	if exact && (typ != "json" && typ != "plaintext") {
		return false, fmt.Errorf(`cannot check exact body with type %q, only %q and %q supported`, typ, typePlain, typeJson)
	}

	return true, nil
}
