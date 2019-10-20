package variables

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type VarsSuite struct {
	suite.Suite
}

func (s *VarsSuite) TestNewFromEnv() {

	testCases := []struct {
		env     map[string]string
		expVars Variables
	}{
		{
			env: map[string]string{
				"ENV1": "val1",
				"env2": "val2",
			},
			expVars: Variables{
				data: map[string]interface{}{
					envNamespace: map[string]interface{}{
						"ENV1": "val1",
						"env2": "val2",
					},
				},
			},
		},
		{
			env: map[string]string{},
			expVars: Variables{
				data: map[string]interface{}{
					envNamespace: make(map[string]interface{}),
				},
			},
		},
	}

	for i, t := range testCases {
		clearEnv()
		setupEnv(t.env)

		actualVars := New(WithEnv())
		s.Equalf(t.expVars, actualVars, "test case %d/%d", i+1, len(testCases))
	}
}

func (s *VarsSuite) TestMerge() {
	arbitraryMap := map[string]interface{}{
		"1": "val1",
		"2": map[string]interface{}{
			"1": 123,
			"2": true,
		},
		"3": []int{4, 5, 6},
	}

	testCases := []struct {
		v1      Variables
		v2      Variables
		expVars Variables
	}{
		{
			v1: newEmptyVars(),
			v2: Variables{
				data: map[string]interface{}{
					varNamespace: arbitraryMap,
					envNamespace: make(map[string]interface{}),
				},
			},
			expVars: Variables{
				data: map[string]interface{}{
					varNamespace: arbitraryMap,
					envNamespace: make(map[string]interface{}),
				},
			},
		},
		{
			v1: Variables{
				data: map[string]interface{}{
					varNamespace: arbitraryMap,
					envNamespace: make(map[string]interface{}),
				},
			},
			v2: newEmptyVars(),
			expVars: Variables{
				data: map[string]interface{}{
					varNamespace: arbitraryMap,
					envNamespace: make(map[string]interface{}),
				},
			},
		},
		{
			v1: Variables{
				data: map[string]interface{}{
					varNamespace:     map[string]interface{}{"1": "2"},
					envNamespace:     make(map[string]interface{}),
					"some-other-key": make(map[string]interface{}),
				},
			},
			v2: Variables{
				data: map[string]interface{}{
					varNamespace: map[string]interface{}{"1": "val1"},
					envNamespace: map[string]interface{}{"env1": "val1"},
				},
			},
			expVars: Variables{
				data: map[string]interface{}{
					"some-other-key": make(map[string]interface{}),
					varNamespace:     map[string]interface{}{"1": "2"},
					envNamespace:     map[string]interface{}{"env1": "val1"},
				},
			},
		},
		{
			v1: Variables{}, // i.e. with nil map
			v2: Variables{
				data: map[string]interface{}{
					varNamespace: map[string]interface{}{"1": "val1"},
					envNamespace: map[string]interface{}{"env1": "val1"},
				},
			},
			expVars: Variables{
				data: map[string]interface{}{
					varNamespace: map[string]interface{}{"1": "val1"},
					envNamespace: map[string]interface{}{"env1": "val1"},
				},
			},
		},
		{
			v1: Variables{
				data: map[string]interface{}{
					varNamespace: map[string]interface{}{"1": "val1"},
					envNamespace: map[string]interface{}{"env1": "val1"},
				},
			},
			v2: Variables{}, // i.e. with nil maps
			expVars: Variables{
				data: map[string]interface{}{
					varNamespace: map[string]interface{}{"1": "val1"},
					envNamespace: map[string]interface{}{"env1": "val1"},
				},
			},
		},
		{
			v1: Variables{
				data: map[string]interface{}{
					"1": map[string]interface{}{
						"22": 2,
					},
					envNamespace: make(map[string]interface{}),
				},
			},
			v2: Variables{
				data: map[string]interface{}{
					"1": map[string]interface{}{
						"a":  'a',
						"22": 5,
					},
					envNamespace: make(map[string]interface{}),
				},
			},
			expVars: Variables{
				data: map[string]interface{}{
					"1": map[string]interface{}{
						"a":  'a',
						"22": 2,
					},
					envNamespace: make(map[string]interface{}),
				},
			},
		},
	}

	for i, t := range testCases {
		actual := t.v1.Merge(t.v2)
		s.Equalf(t.expVars, actual, "test case %d/%d", i+1, len(testCases))
	}
}

func clearEnv() {
	for _, e := range os.Environ() {
		_ = os.Unsetenv(strings.Split(e, "=")[0])
	}
}

func setupEnv(environ map[string]string) {
	for k, v := range environ {
		_ = os.Setenv(k, v)
	}
}

func TestLoaderSuite(t *testing.T) {
	suite.Run(t, new(VarsSuite))
}
