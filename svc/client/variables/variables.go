package variables

// Variables that are used in evaluating transactions and steps
type Variables map[string]interface{}

// Merge another variable instance with this one and return a copy of the result
// not modigfying the original set of variables
func (v Variables) Merge(other Variables) Variables {
	for key, value := range other {
		v[key] = value
	}
	return v
}
