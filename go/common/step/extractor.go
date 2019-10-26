package step

import (
	"encoding/json"

	"github.com/iv-p/apid/common/http"
	"github.com/iv-p/mapaccess"
)

type extractor interface {
	Extract(*http.Response, Export) Exported
}

type bodyExtractor struct {
}

type Exported map[string]interface{}

func NewBodyExtractor() extractor {
	return &bodyExtractor{}
}

func (e *bodyExtractor) Extract(response *http.Response, export Export) Exported {
	exported := make(Exported, len(export))
	var jsonBody interface{}
	err := json.Unmarshal([]byte(response.ReadBody), &jsonBody)
	if err != nil {
		return exported
	}
	var value interface{}
	for exportedKey, bodyKey := range export {
		value, err = mapaccess.Get(jsonBody, bodyKey)
		if err != nil {
			continue
		}
		exported[exportedKey] = value
	}
	return exported
}
