package matcher

import "errors"

var (
	ErrTypeMismatch = errors.New("type_mismatch")
	ErrInterfaceNil = errors.New("interface_nil")
)

func extractString(i interface{}) (string, error) {
	if i == nil {
		return "", ErrInterfaceNil
	}
	switch v := i.(type) {
	case string:
		return v, nil
	default:
		return "", ErrTypeMismatch
	}
}

func optionalString(i interface{}, def string) string {
	if val, err := extractString(i); err == nil {
		return val
	}
	return def
}

func extractFloat(i interface{}) (float64, error) {
	if i == nil {
		return 0, ErrInterfaceNil
	}

	switch v := i.(type) {
	case float64:
		return v, nil
	default:
		return 0, ErrTypeMismatch
	}
}

func optionalFloat(i interface{}, def float64) float64 {
	if val, err := extractFloat(i); err == nil {
		return val
	}
	return def
}

func extractBool(i interface{}) (bool, error) {
	if i == nil {
		return false, ErrInterfaceNil
	}

	switch v := i.(type) {
	case bool:
		return v, nil
	default:
		return false, ErrTypeMismatch
	}
}

func optionalBool(i interface{}, def bool) bool {
	if val, err := extractBool(i); err == nil {
		return val
	}
	return def
}

func extractMap(i interface{}) (map[string]interface{}, error) {
	if i == nil {
		return nil, ErrInterfaceNil
	}

	switch v := i.(type) {
	case map[string]interface{}:
		return v, nil
	default:
		return nil, ErrTypeMismatch
	}
}

func optionalMap(i interface{}, def map[string]interface{}) map[string]interface{} {
	if val, err := extractMap(i); err == nil {
		return val
	}
	return def
}

func extractSlice(i interface{}) ([]interface{}, error) {
	if i == nil {
		return nil, ErrInterfaceNil
	}

	switch v := i.(type) {
	case []interface{}:
		return v, nil
	default:
		return nil, ErrTypeMismatch
	}
}

func optionalSlice(i interface{}, def []interface{}) []interface{} {
	if val, err := extractSlice(i); err == nil {
		return val
	}
	return def
}
