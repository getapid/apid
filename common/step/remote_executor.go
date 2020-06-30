package step

import (
	"bytes"
	"fmt"
	"io/ioutil"
	native "net/http"

	"github.com/getapid/apid-cli/common/remote/endpoint"

	"github.com/getapid/apid-cli/common/http"
	"gopkg.in/yaml.v2"
)

type remoteHTTPExecutor struct {
	key      string
	endpoint string
	c        *native.Client
}

// NewRemoteHTTPExecutor instantiates a new http executor
func NewRemoteHTTPExecutor(client *native.Client, apiKey string, region string) (Executor, error) {
	endpointProvider := endpoint.NewAPIDEndpointProvider()
	endpoint, err := endpointProvider.GetForRegion(region)
	if err != nil {
		return nil, err
	}
	return &remoteHTTPExecutor{
		c:        client,
		endpoint: endpoint,
		key:      apiKey,
	}, nil
}

func (e *remoteHTTPExecutor) Do(request Request) (*http.Response, error) {
	body, err := yaml.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := native.NewRequest("POST", e.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", e.key)
	res, err := e.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == native.StatusTooManyRequests {
		return nil, fmt.Errorf("resource quota exceeded")
	} else if res.StatusCode == native.StatusUnauthorized {
		return nil, fmt.Errorf("invalid api key")
	}

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := &http.Response{}
	err = yaml.Unmarshal(respBody, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
