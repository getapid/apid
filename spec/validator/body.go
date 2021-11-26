package validator

import (
	"fmt"
	"strings"
)

type BodyTextValidator string

func (v *BodyTextValidator) Set(str string) {
	val := BodyTextValidator(str)
	v = &val
}

func (v BodyTextValidator) Validate(str string) (pass []string, fail []string) {
	if len(v) == 0 {
		fail = append(fail, "expected body is empty")
		return
	}

	if strings.EqualFold(string(v), str) {
		pass = append(pass, fmt.Sprintf("body is `%s`", trim(string(v))))
		return
	}
	if testRegex(string(v), str) {
		pass = append(pass, fmt.Sprintf("body matches `%s`", trim(string(v))))
		return
	}

	actual := str
	if len(actual) > 50 {
		actual = str[:50]
	}
	fail = append(fail, fmt.Sprintf("expected body to be `%s`, got `%s`", trim(string(v)), trim(actual)))
	return
}

func trim(txt string) string {
	if len(txt) > 50 {
		return fmt.Sprintf("%s...", txt[:50])
	}
	return txt
}
