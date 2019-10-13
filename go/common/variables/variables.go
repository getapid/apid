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

// NewFromMap returns a new Variables instance that has the provided map set as
// the main variables namespace and an empty environment namespace
func NewFromMap(m map[string]interface{}) Variables {
	return Variables{
		data: m,
		env:  make(map[string]interface{}),
	}
}

// Merge another variable instance with this one and return a copy of the result
// not modifying the original set of variables
func (v Variables) Merge(other Variables) Variables {
	return merge(v, other)
}

// Get returns the main namespace of the variables
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
