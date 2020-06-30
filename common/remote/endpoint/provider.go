package endpoint

import (
	"fmt"
)


type EndpointProvider interface {
	GetForRegion(region string) (string, error)
}

var endpoints = map[string]string{
	"ca-central": "cac",

	"us-east": "use",
	"us-west": "usw",

	"ap-sount": "aps",
	"ap-northeast": "apne",
	"ap-southeast": "apse",

	"eu-west": "euw",
	"eu-central": "euc",
	"eu-north": "eun",

	"sa-east": "sae",
}

type apidEndpointProvider struct {}

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

