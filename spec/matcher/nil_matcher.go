package matcher

type nilMatcher struct{}

func NilMatcher() Matcher {
	return &nilMatcher{}
}

func (m nilMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	return false, nil, nil
}

func (m nilMatcher) String() string {
	return "null"
}
