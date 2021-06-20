package config

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/getapid/cli/common/step"
	"github.com/getapid/cli/common/transaction"
	"github.com/getapid/cli/common/variables"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
)

type LoaderSuite struct {
	suite.Suite
}

// TestLoadDir tests if multiple configs read from disk, merged, then persistent to disk
// and then read back again do not change. I.e. if by writing a merged config to disk
// we lose any data.
func (s *LoaderSuite) TestLoadDir() {
	const cfg1 = `
variables:
  mapping1:
    var1: '123'
    var4: '4'
    var5: '5'
schedule: '0 0 * * Fri'
transactions:
- id: "abc"
  steps:
  - id: '007'
    variables:
      twenty-two: '22'
    request:
      method: T
      endpoint: E
      headers: {HEADER: 2}
      body: B
    expect:
      code: 201
      headers: {HEADER: [2]}
      body:
        - is: content
          subset: true
    export: 
      token: 'response.body.token'
  matrix:
    eleven: [1, 11, 111]
    twenty-one: [21]
`

	const cfg2 = `
variables:
  mapping1:
    var2: 'a'
    var3: '123'
locations: ['us', 'eu']
transactions:
- id: "1234"
  steps:
  - id: 'non-empty'
    variables:
      two: '2'
    request:
      method: t
      endpoint: e
      headers: {HEADER: 1}
      body: b
    export: 
      token: 'response.body.token'
skip_ssl_verify: True
`

	expCfg := Config{
		Variables: newVars(map[string]interface{}{
			"mapping1": map[string]interface{}{
				"var1": "123",
				"var2": "a",
				"var3": "123",
				"var4": "4",
				"var5": "5",
			},
		}),
		Schedule:  "0 0 * * Fri",
		Locations: []string{"us", "eu"},
		Transactions: []transaction.Transaction{
			{
				ID: "abc",
				Steps: []step.Step{
					{
						ID:        "007",
						Variables: newVars(map[string]interface{}{"twenty-two": "22"}),
						Request: step.Request{
							Type:                "T",
							Endpoint:            "E",
							Headers:             step.Headers{"HEADER": {"2"}},
							Body:                "B",
							SkipSSLVerification: pbool(false),
						},
						Response: step.ExpectedResponse{
							Code:    pint(201),
							Headers: &step.Headers{"HEADER": {"2"}},
							Body: []*step.ExpectBody{
								&step.ExpectBody{
									Is:     "content",
									Subset: pbool(true),
								},
							},
						},
						Export: step.Export{"token": "response.body.token"},
					},
				},
				Matrix: &transaction.Matrix{
					M: map[string][]interface{}{
						"eleven":     {1, 11, 111},
						"twenty-one": {21},
					},
				},
			},
			{
				ID: "1234",
				Steps: []step.Step{
					{
						ID:        "non-empty",
						Variables: newVars(map[string]interface{}{"two": "2"}),
						Request: step.Request{
							Type:                "t",
							Endpoint:            "e",
							Headers:             step.Headers{"HEADER": {"1"}},
							Body:                "b",
							SkipSSLVerification: pbool(true),
						},
						Export: step.Export{"token": "response.body.token"},
					},
				},
			},
		},
	}

	dir, err := ioutil.TempDir("", "")
	s.Require().NoError(err)
	defer os.RemoveAll(dir)

	for _, p := range []string{cfg1, cfg2} {
		tmpFile, err := ioutil.TempFile(dir, "*.yaml")
		s.Require().NoError(err)

		_, err = tmpFile.WriteString(p)
		s.Require().NoError(err)

		s.Require().NoError(tmpFile.Close())
	}

	actualCfg, err := Load(dir)
	s.NoError(err)
	s.configsEqual(actualCfg, expCfg)

	serialized, err := yaml.Marshal(actualCfg)
	s.NoError(err)
	f, err := ioutil.TempFile(dir, "*.yaml")
	s.NoError(err)

	n, err := f.Write(serialized)
	s.NoError(err)
	s.Equal(len(serialized), n)

	s.NoError(f.Close())

	readBackCfg, err := Load(f.Name())
	s.NoError(err)
	s.configsEqual(readBackCfg, expCfg)
}

func (s *LoaderSuite) configsEqual(actualCfg, expCfg Config) {
	s.Len(actualCfg.Transactions, len(expCfg.Transactions))

	// checking the equality of transactions like that instead of with assert.ElementsMatch
	// makes the diffs more readable, the result is the same. This way it's easier to see which step is different
	// otherwise, you need to make your way in the gob string of the two lists of transactions.
	expTransactions := mapTransactions(expCfg.Transactions)
	for _, tx := range actualCfg.Transactions {
		expTx := expTransactions[tx.ID]
		for si := range expTx.Steps {
			s.ElementsMatch(expTx.Steps[si].Response.Body, tx.Steps[si].Response.Body)
			expTx.Steps[si].Response.Body = nil
			tx.Steps[si].Response.Body = nil
		}
		s.Equal(expTx.Steps, tx.Steps)
		tx.Steps = nil
		expTx.Steps = nil
		s.Equal(expTx, tx)
	}

	expCfg.Transactions = nil
	actualCfg.Transactions = nil
	s.Equal(expCfg, actualCfg)
}

func mapTransactions(transactions []transaction.Transaction) map[string]transaction.Transaction {
	mapped := make(map[string]transaction.Transaction, len(transactions))
	for _, t := range transactions {
		mapped[t.ID] = t
	}
	return mapped
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

	_, err = invalidYaml.Write([]byte(":"))
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

func (s *LoaderSuite) TestLoadReader() {
	yamlCfg, internalCfg := newConfigPair()
	stringCfg, _ := yaml.Marshal(yamlCfg)

	validFile := bytes.NewBuffer(stringCfg)
	emptyFile := bytes.NewReader(nil)
	invalidYaml := bytes.NewReader([]byte("123"))

	testCases := []struct {
		reader    io.Reader
		expConfig Config
		expErr    bool
	}{
		{
			reader:    validFile,
			expConfig: internalCfg,
			expErr:    false,
		},
		{
			reader:    emptyFile,
			expConfig: Config{},
			expErr:    false,
		},
		{
			reader: invalidYaml,
			expErr: true,
		},
	}

	for i, t := range testCases {
		actualCfg, err := LoadReader(t.reader)

		if t.expErr {
			s.Error(err)
		} else {
			s.NoError(err)
		}

		s.Equalf(t.expConfig, actualCfg, "test case %d/%d", i+1, len(testCases))
	}
}

func (s *LoaderSuite) TestApplySSLFlag_StepUnset() {
	tx := transaction.Transaction{
		Steps: []step.Step{
			{
				Request: step.Request{
					SkipSSLVerification: nil,
				},
			},
		},
	}

	applySSLFlag(true, []transaction.Transaction{tx})

	s.True(*tx.Steps[0].Request.SkipSSLVerification)
}

func (s *LoaderSuite) TestApplySSLFlag_StepSet() {
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

	s.False(*tx.Steps[0].Request.SkipSSLVerification)
}

// newConfigPair returns a config that can be marshalled and written to disk
// and a Config which is supposed to be the equivalent of config in the internal type
func newConfigPair() (config, Config) {
	external := config{
		APIKey:              "",
		SkipSSLVerification: true,
		Schedule:            "? * * * *",
		Locations:           []string{"us", "eu"},
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
							Body: []*step.ExpectBody{
								&step.ExpectBody{
									Is:     "content",
									Subset: pbool(true),
								},
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
		Schedule:     external.Schedule,
		Locations:    external.Locations,
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
