package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type JSONValidator struct {
	Selector *string     `json:"selector"`
	Subset   *bool       `json:"is_subset"`
	Is       interface{} `json:"is"`
}

func (v *JSONValidator) Validate(body []byte) (pass []string, fail []string) {
	message := "body"
	var received interface{}
	if v.Selector != nil {
		val := gjson.GetBytes(body, *v.Selector)
		if !val.Exists() {
			fail = append(fail, fmt.Sprintf("key %s missing in response body", *v.Selector))
			return
		}
		received = val.Value()
		message = *v.Selector
	} else {
		err := json.Unmarshal(body, &received)
		if err != nil {
			fail = append(fail, "body is not valid json")
			return
		}
	}

	var expected = v.Is

	subset := false
	if v.Subset != nil {
		subset = *v.Subset
	}

	if !mapStructsEqual(expected, received, subset) {
		fail = append(fail, fmt.Sprintf("wanted %s for %s, but got %s", prettyPrint(expected), message, prettyPrint(received)))
		return
	}

	pass = append(pass, fmt.Sprintf("%s is %s", message, prettyPrint(expected)))
	return
}

func prettyPrint(i interface{}) string {
	s, _ := json.Marshal(i)
	return string(s)
}

func mapStructsEqual(exp, actual interface{}, subset bool) bool {
	switch exp.(type) {
	case map[string]interface{}:
		expMap := exp.(map[string]interface{})
		actualMap, ok := actual.(map[string]interface{})
		if !ok {
			return false
		}
		if !subset {
			// check if all the keys in the actual map are in the expected map
			for k := range actualMap {
				if _, ok := expMap[k]; !ok {
					return false
				}
			}
		}
		for k, expNested := range expMap {
			if actualNested, ok := actualMap[k]; !ok {
				return false
			} else {
				if !mapStructsEqual(expNested, actualNested, subset) {
					return false
				}
			}
		}
	case []interface{}:
		expSlice := exp.([]interface{})
		actualSlice, ok := actual.([]interface{})
		if !ok {
			return false
		}

		if !subset && len(expSlice) != len(actualSlice) {
			return false
		}

		if len(expSlice) == 0 || !ok {
			return ok
		}

		for _, expVal := range expSlice {
			found := false
			for _, actualVal := range actualSlice {
				if mapStructsEqual(expVal, actualVal, subset) {
					found = true
				}
			}
			if !found {
				return false
			}
		}
	case string:
		switch actual.(type) {
		case string:
			if len(exp.(string)) == 0 {
				return false
			}
			if strings.EqualFold(actual.(string), exp.(string)) {
				return true
			}
			return testRegex(exp.(string), actual.(string))
		default:
			return false
		}

	default:
		return exp == actual
	}
	return true
}
