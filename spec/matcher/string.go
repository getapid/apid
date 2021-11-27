package matcher

import (
	"strings"
)

type StringMatcher string

func (v *StringMatcher) Set(str string) {
	val := StringMatcher(str)
	v = &val
}

func (v StringMatcher) Validate(str string) bool {
	if len(v) == 0 {
		return false
	}
	val := string(v)
	if strings.EqualFold(str, val) {
		return true
	}
	return testRegex(string(v), str)
}
