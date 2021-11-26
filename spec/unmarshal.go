package spec

import (
	"encoding/json"

	"github.com/getapid/apid/log"
)

func Unmarshal(data []byte) (Spec, error) {
	var spec Spec
	err := json.Unmarshal(data, &spec)
	if err != nil {
		log.L.Errorf("invalid spec format: %s", err)
		return Spec{}, err
	}
	return spec, nil
}
