package config

import (
	"fmt"
	"reflect"

	"go.uber.org/multierr"
)

// DefaultValidator does not perform any validations.
type DefaultValidator struct {
}

// Validate method performs validation and returns result and optional error.
func (v DefaultValidator) Validate(val interface{}) (bool, error) {
	return true, nil
}

// StringValidator validates string presence and/or its length.
type StringValidator struct {
	requirements map[string]bool
}

// Validate method performs validation and returns result and optional error.
func (v StringValidator) Validate(val interface{}) (b bool, err error) {
	s := val.(string)
	l := len(s)

	if _, ok := v.requirements["required"]; ok && l == 0 {
		err = multierr.Append(err, fmt.Errorf("cannot be blank"))
	}

	if _, ok := v.requirements["version"]; ok && s != "1" {
		err = multierr.Append(err, fmt.Errorf("supported version: \"1\""))
	}

	return err == nil, err
}

// SliceValidator validates slices and their length
type SliceValidator struct {
	requirements map[string]bool
	unique       []string
}

// Validate method performs validation and returns result and optional error.
func (v SliceValidator) Validate(val interface{}) (b bool, err error) {
	if val == nil {
		if _, ok := v.requirements["required"]; ok {
			return false, fmt.Errorf("should not be empty")
		}
	}

	slice := reflect.ValueOf(val)

	if _, ok := v.requirements["required"]; ok && slice.Len() == 0 {
		return false, fmt.Errorf("should not be empty")
	}

	for _, uniqueField := range v.unique {
		seen := make(map[string]bool)
		for i := 0; i < slice.Len(); i++ {
			valueInQuestion := slice.Index(i).FieldByName(uniqueField).String()
			if _, ok := seen[valueInQuestion]; ok {
				err = multierr.Append(err, fmt.Errorf("should contain unique values for %s. Found multiple %s", uniqueField, valueInQuestion))
			}
			seen[valueInQuestion] = true
		}
	}

	for i := 0; i < slice.Len(); i++ {
		err = validateStruct(slice.Index(i).Interface(), nil)
	}

	return err == nil, err
}

// StructValidator validates recursively the struct
type StructValidator struct {
	requirements map[string]bool
}

// Validate method performs validation and returns result and optional error.
func (v StructValidator) Validate(val interface{}) (b bool, err error) {
	if val == nil {
		if _, ok := v.requirements["required"]; ok {
			return false, fmt.Errorf("should not be empty")
		}
	}

	err = validateStruct(val, nil)
	return err == nil, err
}
