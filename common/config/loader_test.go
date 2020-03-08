package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/getapid/apid-cli/common/step"
	"github.com/getapid/apid-cli/common/transaction"
	"github.com/getapid/apid-cli/common/variables"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type LoaderSuite struct {
	suite.Suite
}

func (s *LoaderSuite) TestLoad() {
	yamlCfg, internalCfg := newConfigPair()
	stringCfg, _ := yaml.Marshal(yamlCfg)

	validFile := s.tempFile()
	defer validFile.Close()
	_, err := validFile.Write(stringCfg)
	s.NoError(err)

	emptyFile := s.tempFile()
	defer emptyFile.Close()

	invalidYaml := s.tempFile()
	defer invalidYaml.Close()

	_, err = invalidYaml.Write([]byte("123"))
	s.NoError(err)

	testCases := []struct {
		path      string
		expConfig Config
		expErr    bool
	}{
		{
			path:   "non-existent/path",
			expErr: true,
		},
		{
			path:      validFile.Name(),
			expConfig: internalCfg,
			expErr:    false,
		},
		{
			path:      emptyFile.Name(),
			expConfig: Config{},
			expErr:    false,
		},
		{
			path:   invalidYaml.Name(),
			expErr: true,
		},
	}

	for i, t := range testCases {
		actualCfg, err := Load(t.path)

		if t.expErr {
			s.Error(err)
		} else {
			s.NoError(err)
		}

		s.Equalf(t.expConfig, actualCfg, "test case %d/%d", i+1, len(testCases))
	}
}

func (s *LoaderSuite) TestApplySSLFlag() {
	tx := transaction.Transaction{
		Steps: []step.Step{
			{
				Request: step.Request{
					SkipSSLVerification: pbool(false),
				},
			},
		},
	}

	applySSLFlag(true, []transaction.Transaction{tx})

	s.True(*tx.Steps[0].Request.SkipSSLVerification)
}

// newConfigPair returns a config that can be marshalled and written to disk
// and a Config which is supposed to be the equivalent of config in the internal type
func newConfigPair() (config, Config) {
	external := config{
		APIKey:              "",
		SkipSSLVerification: true,
		Variables: newVars(map[string]interface{}{
			"key": "some value",
		}),
		Transactions: []transaction.Transaction{
			{
				ID:        "1234",
				Variables: newVars(map[string]interface{}{"one": "1"}),
				Steps: []step.Step{
					{
						ID:        "non-empty",
						Variables: newVars(map[string]interface{}{"two": "2"}),
						Request: step.Request{
							Type:     "t",
							Endpoint: "e",
							Headers:  map[string][]string{},
							Body:     "b",
						},
						Response: step.ExpectedResponse{
							Code:    pint(1),
							Headers: &step.Headers{},
							Body: &step.ExpectBody{
								Type:    pstring("typ"),
								Content: "content",
								Exact:   pbool(true),
							},
						},
						Export: step.Export{},
					},
				},
			},
		},
	}

	txs := make([]transaction.Transaction, len(external.Transactions))
	copy(txs, external.Transactions)
	applySSLFlag(external.SkipSSLVerification, txs)

	internal := Config{
		Version:      external.Version,
		APIKey:       external.APIKey,
		Variables:    external.Variables,
		Transactions: txs,
	}

	return external, internal
}

func pint(i int) *int {
	return &i
}

func pstring(s string) *string {
	return &s
}

func pbool(b bool) *bool {
	return &b
}

func (s *LoaderSuite) tempFile() *os.File {
	f, err := ioutil.TempFile("", "")
	s.Require().NoError(err)
	return f
}

func newVars(m map[string]interface{}) variables.Variables {
	return variables.New(variables.WithVars(m))
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}
