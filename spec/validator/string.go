package validator

import (
	"strings"
)

type StringValidator string

func (v *StringValidator) Set(str string) {
	val := StringValidator(str)
	v = &val
}

func (v StringValidator) Validate(str string) bool {
	if len(v) == 0 {
		return false
	}
	val := string(v)
	if strings.EqualFold(str, val) {
		return true
	}
	return testRegex(string(v), str)
}
