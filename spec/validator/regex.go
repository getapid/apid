package validator

import (
	"regexp"
	"strings"

	"github.com/getapid/apid/log"
)

func testRegex(expected, actual string) bool {
	if !strings.HasSuffix(expected, "$") {
		expected = expected + "$"
	}
	if !strings.HasPrefix(expected, "^") {
		expected = "^" + expected
	}
	r, err := regexp.Compile(expected)
	if err != nil {
		log.L.Infof("failed to compile %s as regex", expected)
		return false
	}

	return r.MatchString(actual)
}
