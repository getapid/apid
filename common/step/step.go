package step

import (
	"fmt"

	"github.com/getapid/apid-cli/common/variables"
	"gopkg.in/yaml.v3"
)

// Step is the data for a single endpoint
type Step struct {
	ID        string              `yaml:"id" json:"id" validate:"required"`
	Variables variables.Variables `yaml:"variables" json:"variables"`
	Request   Request             `yaml:"request" json:"request" validate:"required"`
	Response  ExpectedResponse    `yaml:"expect" json:"expect"`
	Export    Export              `yaml:"export" json:"export"`
}

// Request is a single step request data
type Request struct {
	// Type if the method of the request
	Type                string  `yaml:"method" json:"method" validate:"required"`
	Endpoint            string  `yaml:"endpoint" json:"endpoint" validate:"required"`
	Headers             Headers `yaml:"headers" json:"headers"`
	Body                string  `yaml:"body" json:"body"`
	SkipSSLVerification *bool   `yaml:"skip_ssl_verify" json:"skip_ssl_verify"`
}

type ExpectedResponse struct {
	Code    *int          `yaml:"code" json:"code"`
	Headers *Headers      `yaml:"headers" json:"headers"`
	Body    []*ExpectBody `yaml:"body" json:"body" validate:"expectBody"`
}

type Headers map[string][]string

// UnmarshalYAML unmarshalls yaml in the format
//          headers:
//            header1: 1
//            header1: 4
//            header2: [2,3]
//			  header3: 1
//
// into Headers to result in
// 			map{
// 			  header1: [1,4],
// 			  header2: [2,3],
// 			  header3: [1],
//			}
func (r *Headers) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("headers need to be a mapping")
	}

	m := make(map[string][]string, len(value.Content))

	// in mapping mode the content is alternating keys and values
	for i := 0; i < len(value.Content)-1; i += 2 {
		headerNode := value.Content[i]
		valuesNode := value.Content[i+1]

		var parsedValues []string

		switch valuesNode.Kind {
		case yaml.SequenceNode:
			parsedValues = make([]string, len(valuesNode.Content))
			for i, v := range valuesNode.Content {
				if v.Kind != yaml.ScalarNode {
					return fmt.Errorf("found a non-scalar node as one of the values for the header %q", headerNode.Value)
				}
				parsedValues[i] = v.Value
			}
		case yaml.ScalarNode:
			parsedValues = append(m[headerNode.Value], valuesNode.Value)
		default:
			return fmt.Errorf("found a non-supported kind of yaml node when looking for a header value of %q", headerNode.Value)
		}

		m[headerNode.Value] = parsedValues
	}

	*r = m
	return nil
}

type ExpectBody struct {
	Selector *string `yaml:"selector" json:"selector"`
	KeysOnly *bool   `yaml:"keys_only" json:"keys_only"`
	Subset   *bool   `yaml:"subset" json:"subset"`
	Is       string  `yaml:"is" json:"is"`
}

type Export map[string]string
