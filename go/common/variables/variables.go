package variables

type Variables struct {
	data map[string]interface{}
}

func NewVariables() Variables {
	return Variables{
		data: make(map[string]interface{}),
	}
}

// Merge another variable instance with this one and return a copy of the result
// not modigfying the original set of variables
func (v Variables) Merge(key string, other map[string]interface{}) Variables {
	this := v.data[key]
	if this == nil {
		v.data[key] = other
		return v
	}

	v.data[key] = merge(v.data[key].(map[string]interface{}), other)
	return v
}

func (v Variables) Get() map[string]interface{} {
	return v.data
}

func merge(this, other map[string]interface{}) map[string]interface{} {
	for key, value := range other {
		this[key] = value
	}
	return this
}
