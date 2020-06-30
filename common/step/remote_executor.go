package step

import (
	"github.com/getapid/apid-cli/common/remote/endpoint"
	"bytes"
	"io/ioutil"
	native "net/http"
	"os"

	"github.com/getapid/apid-cli/common/log"

	"github.com/getapid/apid-cli/common/http"
	"gopkg.in/yaml.v2"
)

type remoteHTTPExecutor struct {
	key      string
	endpoint string
	c        *native.Client
}

// NewRemoteHTTPExecutor instantiates a new http executor
func NewRemoteHTTPExecutor(client *native.Client, apiKey string, region string) Executor {
	endpointProvider := endpoint.NewAPIDEndpointProvider()
	endpoint, err := endpointProvider.GetForRegion(region)
	if err != nil {
		log.L.Error(err)
		os.Exit(1)
	}
	return &remoteHTTPExecutor{
		c:        client,
		endpoint: endpoint,
		key:      apiKey,
	}
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

	respBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	res.Body.Close()
	result := &http.Response{}
	err = yaml.Unmarshal(respBody, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
