package validator

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/getapid/apid/spec/matcher"
)

type BodyValidator struct {
	root matcher.Matcher
}

func (m *BodyValidator) UnmarshalJSON(b []byte) error {
	var i interface{}
	err := json.Unmarshal(b, &i)
	if err == nil {
		// Complex matcher
		m.root = matcher.GetMatcher(i)
		return nil
	}

	val := string(b)

	// Try cast to bool
	bo, err := strconv.ParseBool(val)
	if err == nil {
		m.root = matcher.GetMatcher(bo)
		return nil
	}

	// Try cast to number
	f, err := strconv.ParseFloat(val, 64)
	if err == nil {
		m.root = matcher.GetMatcher(f)
		return nil
	}

	// Default to string
	m.root = matcher.GetMatcher(val)

	fmt.Println(m.root)
	return nil
}

func (v BodyValidator) Validate(data []byte) (pass []string, fail []string) {
	var i interface{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		val := string(data)
		i = val

		// Try cast to number
		f, err := strconv.ParseFloat(val, 64)
		if err == nil {
			i = f
		}

		// Try cast to bool
		b, err := strconv.ParseBool(val)
		if err == nil {
			i = b
		}
	}

	_, p, f := v.root.Match(i, "body")
	pass = append(pass, p...)
	fail = append(fail, f...)

	return
}
