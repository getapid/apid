package matcher

import (
	"fmt"
)

type StatusCodeMatcher int64

func (v *StatusCodeMatcher) Set(i int64) {
	val := StatusCodeMatcher(i)
	v = &val
}

func (v StatusCodeMatcher) Validate(value int64) ([]string, []string) {
	if int64(v) != value {
		return nil, []string{fmt.Sprintf("expected status code to be %d, got %d", v, value)}
	}
	return []string{fmt.Sprintf("status code is %d", v)}, nil
}
