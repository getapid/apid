package variables

type Variables struct {
	Data map[string]interface{}
}

func NewVariables() Variables {
	return Variables{
		Data: make(map[string]interface{}),
	}
}

// Merge another variable instance with this one and return a copy of the result
// not modigfying the original set of variables
func (v Variables) Merge(key string, other map[string]interface{}) Variables {
	this := v.Data[key]
	if this == nil {
		v.Data[key] = other
		return v
	}

	v.Data[key] = merge(v.Data[key].(map[string]interface{}), other)
	return v
}

func (v Variables) Get() map[string]interface{} {
	return v.Data
}

func merge(this, other map[string]interface{}) map[string]interface{} {
	for key, value := range other {
		this[key] = value
	}
	return this
}
