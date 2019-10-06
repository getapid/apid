package config

import (
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/multierr"
)

// Validator iterface for the different types to implement
type Validator interface {
	// Validate method performs validation and returns result and optional error.
	Validate(interface{}) (bool, error)
}

// Returns validator struct corresponding to validation type
func getValidatorFromTag(tag string, fieldType reflect.Type) Validator {
	switch t := tag; {
	case t == "string":
		return StringValidator{}
	case t == "required":
		return RequiredValidator{}
	case t == "slice":
		return SliceValidator{}
	case strings.HasPrefix(t, "unique"):
		return UniqueValidator{prop: strings.TrimPrefix(t, "unique=")}
	case t == "recurse":
		return RecurseValidator{}
	case t == "struct":
		return StructValidator{}
	case t == "version":
		return VersionValidator{}
	}

	return DefaultValidator{}
}

// Performs actual data validation using validator definitions on the struct
func validateStruct(s interface{}, accErr error) error {
	// ValueOf returns a Value representing the run-time data
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		// Get the field tag value
		tag := v.Type().Field(i).Tag.Get("validate")
		// Skip if tag is ignored
		if tag == "-" {
			continue
		}

		args := strings.Split(tag, ",")
		for _, arg := range args {
			// Get a validator that corresponds to a tag
			validator := getValidatorFromTag(arg, v.Type().Field(i).Type)
			// Perform validation
			valid, err := validator.Validate(v.Field(i).Interface())
			// Append error to results
			if !valid && err != nil {
				accErr = multierr.Append(accErr, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
			}
		}
	}
	return accErr
}

// Validate deeply validates the config based on the tags
func Validate(c Config) error {
	return validateStruct(c, nil)
}
