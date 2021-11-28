package matcher

type anyMatcher struct{}

func AnyMatcher() Matcher {
	return &anyMatcher{}
}

func (m anyMatcher) Match(data interface{}, location string) (bool, []string, []string) {
	return true, nil, nil
}

func (m anyMatcher) String() string {
	return "any"
}
