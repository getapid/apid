package variables

import (
	"reflect"
)

type Variables map[string]interface{}

func New() Variables {
	return make(map[string]interface{})
}

// Merge another variable instance with this one and return a copy of the result
// not modifying the original set of variables
func (v Variables) Merge(other Variables) Variables {
	return mergeMaps(v, other)
}

func mergeMaps(this, other map[string]interface{}) map[string]interface{} {
	if this == nil {
		return other
	}
	for key, newVal := range other {
		if existingVal, ok := this[key]; ok {
			// if the existing value isn't mergable we skip it
			if existingMap, ok := tryToStringMap(existingVal); ok {
				// if the new value isn't mergable we overwrite the existing value
				if newMap, ok := tryToStringMap(newVal); ok {
					this[key] = mergeMaps(existingMap, newMap)
					continue
				}
			}
		}
		this[key] = deepCopyInterface(newVal)
	}

	return this
}

// tryToStringMap returns val as map[string]interface{} and true if val is easily convertible to that.
// It returns nil and false if it is not easily convertible.
//
// If val is already map[string]interface{}, val and true is directly returned.
// If val is map[interface{}]interface{} and each key is actually a string,
// then a map[string]interface{} will be returned. Otherwise nil and false will be returned.
func tryToStringMap(val interface{}) (map[string]interface{}, bool) {
	stringMap, ok := val.(map[string]interface{})
	if ok {
		return stringMap, true
	}

	interfaceMap, ok := val.(map[interface{}]interface{})
	if !ok {
		return nil, false
	}
	stringMap = make(map[string]interface{}, len(interfaceMap))
	for k, v := range interfaceMap {
		if strK, ok := k.(string); ok {
			stringMap[strK] = v
		} else {
			return nil, false
		}
	}
	return stringMap, true
}

func deepCopyInterface(i interface{}) interface{} {
	return deepCopy(reflect.ValueOf(i)).Interface()
}

func deepCopy(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Interface, reflect.Ptr:
		if value.IsNil() {
			return reflect.ValueOf(nil)
		}
		return deepCopy(value.Elem())
	case reflect.Map:
		copied := reflect.MakeMapWithSize(value.Type(), value.Len())

		iter := value.MapRange()
		for iter.Next() {
			k, v := iter.Key(), iter.Value()
			copied.SetMapIndex(k, deepCopy(v))
		}
		return copied
	case reflect.Slice, reflect.Array:
		copyLen := value.Len()
		copied := reflect.MakeSlice(value.Type(), copyLen, copyLen)

		for i := 0; i < copyLen; i++ {
			copiedElement := deepCopy(value.Index(i))
			copied.Index(i).Set(copiedElement)
		}
		return copied
	default:
		return value
	}
}
