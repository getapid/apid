package validator

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getapid/apid/log"
	"github.com/getapid/apid/spec/matcher"
)

type HeaderValidator struct {
	root map[matcher.Matcher]matcher.Matcher
}

func (m *HeaderValidator) UnmarshalJSON(data []byte) (err error) {
	var i interface{}
	err = json.Unmarshal(data, &i)
	if err != nil {
		log.L.Fatal("invalid header validatior, got %v", data)
		return
	}

	m.root = make(map[matcher.Matcher]matcher.Matcher)
	switch v := i.(type) {
	case map[string]interface{}:
		for name, value := range v {
			m.root[matcher.GetMatcher(name)] = matcher.GetMatcher(value)
		}
	default:
		log.L.Fatal("invalid header validatior, got %v", v)
	}

	return
}

func (v HeaderValidator) Validate(headers http.Header) (pass []string, fail []string) {
	log.L.Infof("validating headers %v", headers)
	for nameMatcher, valueMatcher := range v.root {
		found := false
		errStr := fmt.Sprintf("header %s not found", nameMatcher)
		for name, values := range headers {
			if ok, _, _ := nameMatcher.Match(name, "header"); !ok {
				continue
			}
			errStr = fmt.Sprintf("header %s does not match %s, got %v", nameMatcher, valueMatcher, values)

			for _, value := range values {
				if ok, _, _ := valueMatcher.Match(value, "header"); !ok {
					continue
				}
				pass = append(pass, fmt.Sprintf("header %s is %s", nameMatcher, valueMatcher))
				found = true
				break
			}
			if found {
				break
			}
		}
		if !found {
			fail = append(fail, errStr)
		}
	}
	return pass, fail
}
