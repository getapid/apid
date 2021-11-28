package matcher

type boolMatcher struct {
	value bool
}

func BoolMatcher(value bool) Matcher {
	return boolMatcher{value}
}

func (m boolMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	// switch val := data.(type) {
	// case bool:
	// 	return v.value == val, ""
	// }
	// return false, ""
	return false, nil, nil
}

func (m boolMatcher) String() string {
	return "false"
}
