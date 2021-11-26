package step

import (
	"encoding/json"
	"strings"

	"github.com/getapid/apid/http"
	"github.com/getapid/apid/log"
	"github.com/getapid/apid/spec"
	"github.com/getapid/apid/variables"
	"github.com/tidwall/gjson"
)

type Exporter struct{}

func NewExporter() Exporter {
	return Exporter{}
}

func (e *Exporter) export(response *http.Response, export spec.Export) variables.Variables {
	exported := make(variables.Variables, len(export))
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
