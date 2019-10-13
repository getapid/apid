package template

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iv-p/apid/common/variables"
	"github.com/iv-p/mapaccess"
)

// Render parses the string and returns the interpolated result
func Render(template string, data variables.Variables) (string, error) {
	var res strings.Builder
	parser := parse(template, leftDelim, rightDelim)
	for {
		token := parser.nextItem()
		switch token.typ {
		case tokenError:
			return res.String(), errors.New(token.val)
		case tokenEnd:
			return res.String(), nil
		case tokenText:
			if _, err := res.WriteString(token.val); err != nil {
				return res.String(), err
			}
		case tokenIdentifier:
			var (
				dataSource = data.Get()
				tokenVal   = token.val
			)

			switch t := token.val; {
			case strings.HasPrefix(t, "variables."):
				tokenVal = strings.TrimPrefix(t, "variables.")
				dataSource = data.Get()
			case strings.HasPrefix(t, "env."):
				tokenVal = strings.TrimPrefix(t, "env.")
				dataSource = data.GetEnv()
			}

			val, err := mapaccess.Get(dataSource, tokenVal)
			if err != nil {
				return res.String(), err
			}
			switch c := val.(type) {
			case string:
				res.WriteString(c)
			default:
				return res.String(), fmt.Errorf("value for key %s is not string", token.val)
			}
		}
	}
}
