package endpoint

import (
	"fmt"
)


type EndpointProvider interface {
	GetForRegion(region string) (string, error)
}

var endpoints = map[string]string{
	"ca-central": "cac", // done ca-central-1

	"us-east": "use", // done us-east-1
	"us-west": "usw", // done us-west-1

	"ap-sount": "aps", // done ap-sount-1
	"ap-northeast": "apne", // done ap-northeast-1
	"ap-southeast": "apse", // done ap-southeast-2

	"eu-west": "euw", // done eu-west-1
	"eu-south": "eus", // pricing eu-south-1
	"eu-north": "eun", // done eu-north-1

	"sa-east": "sae", // done sa-east-1
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

