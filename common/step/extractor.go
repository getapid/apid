package step

import (
	"encoding/json"
	"strings"

	"github.com/getapid/cli/common/http"
	"github.com/getapid/cli/common/log"
	"github.com/tidwall/gjson"
)

type extractor interface {
	extract(*http.Response, Export) Exported
}

type bodyExtractor struct {
}

type Exported map[string]interface{}

func (e Exported) Generic() map[string]interface{} {
	return e
}

func NewBodyExtractor() extractor {
	return &bodyExtractor{}
}

// extract will try to extract all the variables specified in the provided Export.
// It will ignore any keys it cannot find in the response (headers, body or others).
func (e *bodyExtractor) extract(response *http.Response, export Export) Exported {
	exported := make(Exported, len(export))
	var jsonBody interface{}
	err := json.Unmarshal(response.Body, &jsonBody)
	if err != nil {
		return exported
	}
	for exportAs, keyToExport := range export {
		keyToExport = strings.TrimPrefix(keyToExport, "response.")

		switch {
		case strings.HasPrefix(keyToExport, "body."):
			keyToExport = strings.TrimPrefix(keyToExport, "body.")
			val := gjson.Get(string(response.Body), keyToExport)
			if !val.Exists() {
				log.L.Warnf("could not find key %v : %v", keyToExport, err)
				continue
			}
			exported[exportAs] = val.String()
		case strings.HasPrefix(keyToExport, "headers."):
			keyToExport = strings.TrimPrefix(keyToExport, "headers.")
			foundHeaders, ok := response.Header[keyToExport]
			if !ok {
				log.L.Warnf("could not find key %v from headers %v", keyToExport, response.Header)
				continue
			}
			if len(foundHeaders) > 1 {
				log.L.Warnf("found multiple header values for key %v", keyToExport)
				continue
			}
			exported[exportAs] = foundHeaders[0]
		default:
			continue
		}
	}
	return exported
}
