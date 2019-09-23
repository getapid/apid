package interpolator

import (
	"log"
	"regexp"
	"strings"

	"github.com/iv-p/apiping/svc/cli/json"
)

type StringInterpolator interface {
	Interpolate(string, interface{}) (string, error)
}

type SimpleStringInterpolator struct {
	re       *regexp.Regexp
	accessor json.Accessor

	StringInterpolator
}

func NewSimpleStringInterpolator() StringInterpolator {
	return &SimpleStringInterpolator{
		re:       regexp.MustCompile(`\${(.+?)}`),
		accessor: json.NewJsonAccessor(),
	}
}

func (s *SimpleStringInterpolator) Interpolate(str string, variables interface{}) (string, error) {
	log.Println(variables)
	matched := s.re.FindAllStringSubmatch(str, -1)
	if len(matched) == 0 {
		return str, nil
	}

	for _, match := range matched {
		toReplace := match[0]
		toEvaluate := strings.Trim(match[1], " ")

		data, err := s.accessor.Get(toEvaluate, variables)
		if err != nil {
			log.Println(err)
			str = strings.Replace(str, toReplace, "", 1)
			continue
		}

		str = strings.Replace(str, toReplace, data, 1)
	}

	return str, nil
}
