package transaction

import (
	"testing"

	"github.com/getapid/apid-cli/common/variables"
	"github.com/stretchr/testify/assert"
)

func TestMatrix(t *testing.T) {
	testCases := [...]struct {
		variables map[string][]interface{}
		expSets   []variables.Variables
	}{
		{
			variables: map[string][]interface{}{
				"a": {1, 2},
				"b": {"b1", "b2"},
			},
			expSets: []variables.Variables{
				newVars(map[string]interface{}{
					"a": 1,
					"b": "b1",
				}),
				newVars(map[string]interface{}{
					"a": 1,
					"b": "b2",
				}),
				newVars(map[string]interface{}{
					"a": 2,
					"b": "b1",
				}),
				newVars(map[string]interface{}{
					"a": 2,
					"b": "b2",
				}),
			},
		},
		{
			variables: map[string][]interface{}{},
			expSets:   []variables.Variables{},
		},
		{
			variables: map[string][]interface{}{
				"a": {1},
				"b": {"1"},
			},
			expSets: []variables.Variables{
				newVars(map[string]interface{}{
					"a": 1,
					"b": "1",
				}),
			},
		},
	}

	for i, tc := range testCases {
		m := Matrix{M: tc.variables}

		returnedSets := make([]variables.Variables, 0, len(tc.expSets))
		for m.HasNext() {
			returnedSets = append(returnedSets, m.NextSet())
		}
		assert.ElementsMatch(t, tc.expSets, returnedSets, "test case %d/%d", i+1, len(testCases))
	}
}

func newVars(m map[string]interface{}) variables.Variables {
	return variables.New(variables.WithVars(m))
}
