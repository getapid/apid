package validator

import (
	"encoding/json"
	"strconv"

	"github.com/getapid/apid/log"
	"github.com/getapid/apid/spec/matcher"
)

type StatusCodeValidator struct {
	root matcher.Matcher
}

func (m *StatusCodeValidator) UnmarshalJSON(data []byte) (err error) {
	var i interface{}
	err = json.Unmarshal(data, &i)
	if err == nil {
		m.root = matcher.GetMatcher(i)
		return
	}

	// Try cast to number
	d, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		m.root = matcher.IntMatcherWithOptions(d)
		return
	}

	log.L.Fatal("expected number matcher for status code, got %v", string(data))
	return
}

func (v StatusCodeValidator) Validate(code int) ([]string, []string) {
	_, p, f := v.root.Match(code, "status code")
	return p, f
}
