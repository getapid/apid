package validator

import (
	"fmt"
)

type StatusCodeValidator int64

func (v *StatusCodeValidator) Set(i int64) {
	val := StatusCodeValidator(i)
	v = &val
}

func (v StatusCodeValidator) Validate(value int64) ([]string, []string) {
	if int64(v) != value {
		return nil, []string{fmt.Sprintf("expected status code to be %d, got %d", v, value)}
	}
	return []string{fmt.Sprintf("status code is %d", v)}, nil
}
