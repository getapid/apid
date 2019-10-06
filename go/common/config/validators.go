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
type StringValidator struct{}

// Validate method performs validation and returns result and optional error.
func (v StringValidator) Validate(val interface{}) (b bool, err error) {
	if val == nil {
		return true, nil
	}

	_, ok := val.(string)

	if !ok {
		return false, fmt.Errorf("must be a string")
	}

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

// RecurseValidator calls the validation on each element of a slice
type RecurseValidator struct{}

// Validate method performs validation and returns result and optional error.
func (v RecurseValidator) Validate(val interface{}) (bool, error) {
	if val == nil {
		return true, nil
	}

	slice := reflect.ValueOf(val)

	if slice.Kind() != reflect.Slice {
		return false, fmt.Errorf("must be a slice")
	}

	var err error
	for i := 0; i < slice.Len(); i++ {
		err = validateStruct(slice.Index(i).Interface(), err)
	}

	return err == nil, err
}

// SliceValidator validates slices
type SliceValidator struct{}

// Validate method performs validation and returns result and optional error.
func (v SliceValidator) Validate(val interface{}) (b bool, err error) {
	if val == nil {
		return true, nil
	}

	slice := reflect.ValueOf(val)

	if slice.Kind() != reflect.Slice {
		return false, fmt.Errorf("must be a slice")
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

// StructValidator validates recursively the struct
type StructValidator struct{}

// Validate method performs validation and returns result and optional error.
func (v StructValidator) Validate(val interface{}) (bool, error) {
	if val == nil {
		return true, nil
	}

	struc := reflect.ValueOf(val)

	if struc.Kind() != reflect.Struct {
		return false, fmt.Errorf("must be a struct")
	}

	err := validateStruct(val, nil)
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
