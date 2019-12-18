package template

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/getapid/apid-cli/common/command"

	"github.com/getapid/apid-cli/common/variables"
	"github.com/iv-p/mapaccess"
	"go.uber.org/multierr"
)

// Render parses the string and returns the interpolated result
func Render(template string, data variables.Variables) (string, error) {
	cmd := command.NewShellExecutor()
	var renderer strings.Builder
	var multiErr error
	parser := parse(template)
	for {
		token := parser.nextItem()
		switch token.typ {
		case tokenError:
			multiErr = multierr.Append(multiErr, fmt.Errorf("parsing error: %v : %v", template, token.val))
		case tokenEnd:
			goto EXIT
		case tokenText:
			if _, err := renderer.WriteString(token.val); err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf("write string: %v : %v", template, err))
			}
		case tokenIdentifier:
			val, err := mapaccess.Get(data.Raw(), token.val)
			if err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf("%v: %v", token.val, err))
				continue
			}
			switch c := val.(type) {
			case string:
				renderer.WriteString(c)
			case float64:
				renderer.WriteString(fmt.Sprintf("%g", c))
			case int:
				renderer.WriteString(fmt.Sprintf("%d", c))
			default:
				multiErr = multierr.Append(multiErr, fmt.Errorf("unknown value type %v: %v", reflect.TypeOf(val), token.val))
			}
		case tokenCommand:
			res, err := cmd.Exec(token.val, data)
			if err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf("error executing command %v: %v", token.val, err))
				continue
			}
			renderer.Write(res)
		}
	}
EXIT:
	return renderer.String(), multiErr
}
