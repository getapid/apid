package template

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iv-p/apid/common/variables"
	"github.com/iv-p/mapaccess"
	"go.uber.org/multierr"
)

// Render parses the string and returns the interpolated result
func Render(template string, data variables.Variables) (string, error) {
	var res strings.Builder
	var multiErr error
	parser := parse(template, leftDelim, rightDelim)
	for {
		token := parser.nextItem()
		switch token.typ {
		case tokenError:
			multiErr = multierr.Append(multiErr, fmt.Errorf("parsing error: %v : %v", template, token.val))
		case tokenEnd:
			goto EXIT
		case tokenText:
			if _, err := res.WriteString(token.val); err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf("write string: %v : %v", template, err))
			}
		case tokenIdentifier:
			val, err := mapaccess.Get(data.Raw(), token.val)
			if err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf("key not found in data: %v : %v", token.val, err))
				continue
			}
			switch c := val.(type) {
			case string:
				res.WriteString(c)
			case float64:
				res.WriteString(fmt.Sprintf("%g", c))
			case int:
				res.WriteString(fmt.Sprintf("%d", c))
			default:
				multiErr = multierr.Append(multiErr, fmt.Errorf("unknown value type %v: %v", reflect.TypeOf(val), token.val))
			}
		}
	}
EXIT:
	return res.String(), multiErr
}
