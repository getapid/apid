package step

import (
	"encoding/json"

	"github.com/iv-p/apid/common/http"
	"github.com/iv-p/apid/common/log"
	"github.com/iv-p/mapaccess"
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

func (e *bodyExtractor) extract(response *http.Response, export Export) Exported {
	exported := make(Exported, len(export))
	var jsonBody interface{}
	err := json.Unmarshal(response.ReadBody, &jsonBody)
	if err != nil {
		return exported
	}
	var value interface{}
	for exportedKey, bodyKey := range export {
		value, err = mapaccess.Get(jsonBody, bodyKey)
		if err != nil {
			log.L.Warn("could not fetch key %v : %v", bodyKey, err)
			continue
		}
		exported[exportedKey] = value
	}
	return exported
}
