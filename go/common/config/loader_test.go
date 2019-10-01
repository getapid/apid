package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/iv-p/apid/common/step"
	"github.com/iv-p/apid/common/transaction"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type LoaderSuite struct {
	suite.Suite
}

func (s *LoaderSuite) TestLoad() {
	cfg := newConfig()
	stringCfg, _ := yaml.Marshal(cfg)

	validFile := s.tempFile()
	_, err := validFile.Write(stringCfg)
	s.NoError(err)

	emptyFile := s.tempFile()

	invalidYaml := s.tempFile()
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
			expConfig: cfg,
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

	for _, t := range testCases {
		actualCfg, err := Load(t.path)

		if t.expErr {
			s.Error(err)
		} else {
			s.NoError(err)
		}

		s.Equal(t.expConfig, actualCfg)
	}
}

func newConfig() Config {
	return Config{
		APIKey: "",
		Variables: map[string]interface{}{
			"key": "some value",
		},
		Transactions: []transaction.Transaction{
			{
				ID:        "1234",
				Variables: map[string]interface{}{"one": "1"},
				Steps: []step.Step{
					{
						ID:        "non-empty",
						Variables: map[string]interface{}{"two": "2"},
						Request: step.Request{
							Type:     "t",
							Endpoint: "e",
							Headers:  map[string]string{},
							Body:     "b",
						},
						Response: step.ExpectedResponse{
							Code:    pint(1),
							Headers: &step.Headers{},
							Body: &step.ExpectBody{
								Type:    pstring("typ"),
								Content: pstring("content"),
								Exact:   pbool(true),
							},
						},
					},
				},
			},
		},
	}
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
	f, err := ioutil.TempFile("", "*****")
	s.Require().NoError(err)
	return f
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(LoaderSuite))
}
