package template

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/getapid/apid/log"
	"github.com/getapid/apid/variables"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"
)

// Render parses the string and returns the interpolated result
func Render(template string, data variables.Variables) (string, error) {
	var renderer strings.Builder
	var multiErr error
	parser := parse(template)
	for {
		token := parser.nextItem()
		switch token.typ {
		case tokenError:
			multiErr = multierr.Append(multiErr, fmt.Errorf("parsing error: %v : %v", template, token.val))
		case tokenEnd:
			return renderer.String(), multiErr
		case tokenText:
			if _, err := renderer.WriteString(token.val); err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf("write string: %v : %v", template, err))
			}
		case tokenIdentifier:
			d, err := json.Marshal(data)
			if err != nil {
				log.L.Error("could not serialize variables")
				continue
			}
			val := gjson.GetBytes(d, token.val)
			if !val.Exists() {
				multiErr = multierr.Append(multiErr, fmt.Errorf("%v: key not found", token.val))
				continue
			}
			renderer.WriteString(val.String())
		}
	}
}
