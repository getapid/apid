package variables

// TODO test this package

type Variables struct {
	data map[string]interface{}
}

// TODO rename to just New()
func NewVariables() Variables {
	return Variables{
		data: make(map[string]interface{}),
	}
}

// TODO rename to just NewWithMap()
func NewVariablesFromMap(m map[string]interface{}) Variables {
	return Variables{
		data: m,
	}
}

// Merge another variable instance with this one and return a copy of the result
// not modifying the original set of variables
func (v Variables) Merge(other Variables) Variables {
	return merge(v, other)
}

// TODO rename to GetData()
func (v Variables) Get() map[string]interface{} {
	return v.data
}

// TODO make this recursive
func merge(this, other Variables) Variables {
	for key, value := range other.data {
		this.data[key] = value
	}
	return this
}
