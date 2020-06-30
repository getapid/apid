package template

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getapid/apid-cli/common/command"
	"github.com/getapid/apid-cli/common/log"
	"github.com/tidwall/gjson"

	"github.com/getapid/apid-cli/common/variables"
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
			d, err := json.Marshal(data.Raw())
			if err != nil {
				log.L.Error("could not serialize variables")
				continue
			}
			val := gjson.Get(string(d), token.val)
			renderer.WriteString(val.String())
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
