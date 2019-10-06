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
	case t == "required":
		return RequiredValidator{}
	case strings.HasPrefix(t, "unique"):
		return UniqueValidator{prop: strings.TrimPrefix(t, "unique=")}
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
		typeField := v.Type().Field(i)
		typeOfTypeField := typeField.Type
		field := v.Field(i).Interface()

		// Get the field tag value
		tag := typeField.Tag.Get("validate")
		// Skip if tag is ignored
		if tag == "-" {
			continue
		}

		args := strings.Split(tag, ",")
		for _, arg := range args {
			// Get a validator that corresponds to a tag
			validator := getValidatorFromTag(arg, typeOfTypeField)
			// Perform validation
			valid, err := validator.Validate(field)
			// Append error to results
			if !valid && err != nil {
				accErr = multierr.Append(accErr, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
			}
		}

		if typeOfTypeField.Kind() == reflect.Struct {
			accErr = validateStruct(field, accErr)
		} else if typeOfTypeField.Kind() == reflect.Slice {
			f := v.Field(i)
			for i := 0; i < f.Len(); i++ {
				accErr = validateStruct(f.Index(i).Interface(), accErr)
			}
		}
	}
	return accErr
}

// Validate deeply validates the config based on the tags
func Validate(c Config) error {
	return validateStruct(c, nil)
}
