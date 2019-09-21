package json

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Accessor interface {
	Get(string, interface{}) (string, error)
}

type JsonAccessor struct {
	Accessor
}

func NewJsonAccessor() Accessor {
	return &JsonAccessor{}
}

func (a *JsonAccessor) Get(key string, data interface{}) (string, error) {
	keys := strings.Split(key, ".")
	d := data
	var err error
	for _, k := range keys {
		d, err = getValueForKey(k, d)
		if err != nil {
			return "", err
		}
	}

	str, err := parseToString(d)
	if err != nil {
		return "", err
	}
	log.Println(str)
	return str, nil
}

func getValueForKey(key string, data interface{}) (interface{}, error) {
	key, index, isArray := getKeyAndIndex(key)
	switch v := data.(type) {
	case map[string]interface{}:
		res, ok := v[key]
		if !ok {
			log.Println(v)
			return res, fmt.Errorf("key '%s' not found", key)
		}
		if isArray {
			return getIndexFromArray(index, res)
		}
		return res, nil
	case []interface{}:
		if index < 0 || index > len(v) {
			return data, fmt.Errorf("index out of bounds %d", index)
		}
		return v[index], nil
	default:
		return data, fmt.Errorf("can't deal with this type %s", reflect.TypeOf(data))
	}
}

func getIndexFromArray(index int, data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case []interface{}:
		if index < 0 || index > len(v) {
			return data, fmt.Errorf("index out of bounds %d", index)
		}
		return v[index], nil
	default:
		return data, fmt.Errorf("can't deal with this type %s", reflect.TypeOf(data))
	}
}

func getKeyAndIndex(key string) (string, int, bool) {
	re := regexp.MustCompile(`^\w*\[\d+\]$`)
	isArray := re.MatchString(key)
	if !isArray {
		return key, -1, false
	}

	re = regexp.MustCompile(`^(\w*)\[(\d+)\]$`)
	regexpResult := re.FindAllStringSubmatch(key, -1)
	if len(regexpResult) != 1 && len(regexpResult[0]) != 3 {
		return key, -1, false
	}

	k := regexpResult[0][1]
	i, err := strconv.Atoi(regexpResult[0][2])
	if err != nil {
		return key, -1, false
	}
	return k, i, true
}

func parseToString(data interface{}) (string, error) {
	switch v := data.(type) {
	case float64:
		if float64(int(v)) == v {
			return fmt.Sprintf("%d", int(v)), nil
		}
		return fmt.Sprintf("%g", v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("could not transform value to string, unsupported type %s", reflect.TypeOf(v))
	}
}
