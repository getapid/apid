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
	args := strings.Split(tag, ",")

	requirements := make(map[string]bool)
	for _, req := range args {
		requirements[req] = true
	}

	switch fieldType.Kind() {
	case reflect.Slice:
		validator := SliceValidator{requirements: requirements, unique: make([]string, 0)}

		for _, req := range args {
			if strings.HasPrefix(req, "unique") {
				validator.unique = append(validator.unique, strings.TrimPrefix(req, "unique="))
			}
		}

		return validator

	case reflect.String:
		return StringValidator{requirements: requirements}

	case reflect.Struct:
		return StructValidator{requirements: requirements}
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
		// Get a validator that corresponds to a tag
		validator := getValidatorFromTag(tag, v.Type().Field(i).Type)
		// Perform validation
		valid, err := validator.Validate(v.Field(i).Interface())
		// Append error to results
		if !valid && err != nil {
			accErr = multierr.Append(accErr, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}
	return accErr
}

// Validate deeply validates the config based on the tags
func Validate(c Config) error {
	return validateStruct(c, nil)
}
