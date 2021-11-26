package validator

import (
	"fmt"

	"github.com/getapid/apid/log"
)

type HeaderValidator map[StringValidator]StringValidator

func (v *HeaderValidator) Set(headers map[string]string) {
	m := make(map[StringValidator]StringValidator, len(headers))
	for name, value := range headers {
		m[StringValidator(name)] = StringValidator(value)
	}
	val := HeaderValidator(m)
	v = &val
}

func (v HeaderValidator) Validate(headers map[string][]string) (pass []string, fail []string) {
	log.L.Infof("validating headers %v", headers)
	for nameValidator, valueValidator := range v {
		found := false
		errStr := fmt.Sprintf("expected header %s, but none found", nameValidator)
		for name, values := range headers {
			if !nameValidator.Validate(name) {
				continue
			}
			errStr = fmt.Sprintf("expected header %s to match %s, got %v", nameValidator, valueValidator, values)

			for _, value := range values {
				if !valueValidator.Validate(value) {
					continue
				}
				pass = append(pass, fmt.Sprintf("has header %s = %s", nameValidator, valueValidator))
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
