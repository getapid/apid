package endpoint

import (
	"fmt"
)

type EndpointProvider interface {
	GetForRegion(region string) (string, error)
}

var endpoints = map[string]string{
	"montreal": "ymq",

	"washington":   "was",
	"sanfrancisco": "sfo",

	"mumbai": "bom",
	"tokyo":  "tyo",
	"sydney": "syd",

	"stockholm": "sto",
	"frankfurt": "fra",
	"dublin":    "dub",

	"saopaulo": "sao",
}

type apidEndpointProvider struct{}

func NewAPIDEndpointProvider() EndpointProvider {
	return &apidEndpointProvider{}
}

func (p apidEndpointProvider) GetForRegion(region string) (string, error) {
	token, ok := endpoints[region]
	if !ok {
		return token, fmt.Errorf("invalid region %s", region)
	}

	return fmt.Sprintf("https://%s.api.getapid.com/executor", token), nil
}
