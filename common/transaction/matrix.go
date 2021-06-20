package transaction

import (
	"encoding/json"

	"github.com/getapid/cli/common/variables"
)

// Matrix is holding slices of values for a number of (step/transaction) variables.
// It generates combinations of these values. It is not thread-safe.
type Matrix struct {
	inited          bool                     `json:"-"`
	order           []string                 `json:"-"`
	currentVariants map[string]int           `json:"-"`
	M               map[string][]interface{} `yaml:",inline"`
}

func (m *Matrix) UnmarshalJSON(b []byte) error {
	m.M = make(map[string][]interface{})
	return json.Unmarshal(b, &m.M)
}

func (m Matrix) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.M)
}

func (m *Matrix) HasNext() bool {
	if !m.inited && len(m.M) > 0 {
		return true
	}
	for _, v := range m.order {
		if m.currentVariants[v]+1 < len(m.M[v]) {
			return true
		}
	}
	return false
}

func (m *Matrix) NextSet() variables.Variables {
	if !m.inited {
		m.currentVariants = make(map[string]int, len(m.M))
		m.order = make([]string, 0, len(m.M))
		for varName := range m.M {
			m.order = append(m.order, varName)
		}
		m.inited = true
		return m.CurrentSet()
	}

	for i := len(m.order) - 1; i >= 0; i-- {
		currentVar := m.order[i]
		currentVarVariant := m.currentVariants[currentVar]

		if currentVarVariant+1 >= len(m.M[currentVar]) {
			m.currentVariants[currentVar] = 0
		} else {
			m.currentVariants[currentVar]++
			break
		}
	}

	return m.CurrentSet()
}

func (m Matrix) CurrentSet() variables.Variables {
	varSet := make(map[string]interface{}, len(m.M))
	for _, varName := range m.order {
		varSet[varName] = m.M[varName][m.currentVariants[varName]]
	}
	return variables.New(variables.WithVars(varSet))
}
