package config

import (
	"testing"

	"github.com/getapid/cli/common/step"
	"github.com/stretchr/testify/suite"
)

type ValidatorsSuite struct {
	suite.Suite
}

func (s *ValidatorsSuite) TestValidate_ExpectBodyValidator() {
	testCases := [...]struct {
		name     string
		val      interface{}
		expValid bool
	}{
		{
			name:     "nil",
			val:      nil,
			expValid: true,
		},
		{
			name:     "invalid - no clauses",
			val:      []*step.ExpectBody{},
			expValid: true,
		},
		{
			name: "valid - only is",
			val: []*step.ExpectBody{
				&step.ExpectBody{
					Is: "test",
				},
			},
			expValid: true,
		},
		{
			name: "invalid - empty is clause",
			val: []*step.ExpectBody{
				&step.ExpectBody{
					Is: "",
				},
			},
			expValid: false,
		},
		{
			name: "invalid - empty clause",
			val: []*step.ExpectBody{
				&step.ExpectBody{},
			},
			expValid: false,
		},
		{
			name: "valid - subset true",
			val: []*step.ExpectBody{
				&step.ExpectBody{
					Subset: pbool(true),
					Is:     "test",
				},
			},
			expValid: true,
		},
		{
			name: "valid - subset false",
			val: []*step.ExpectBody{
				&step.ExpectBody{
					Subset: pbool(false),
					Is:     "test",
				},
			},
			expValid: true,
		},
		{
			name: "valid - keys only false",
			val: []*step.ExpectBody{
				&step.ExpectBody{
					KeysOnly: pbool(false),
					Is:       "test",
				},
			},
			expValid: true,
		},
		{
			name: "valid - keys only true",
			val: []*step.ExpectBody{
				&step.ExpectBody{
					KeysOnly: pbool(true),
					Is:       "test",
				},
			},
			expValid: true,
		},
		{
			name: "valid - keys only and subset true",
			val: []*step.ExpectBody{
				&step.ExpectBody{
					KeysOnly: pbool(true),
					Subset:   pbool(true),
					Is:       "test",
				},
			},
			expValid: true,
		},
	}

	validator := ExpectBodyValidator{}
	for _, t := range testCases {
		actualIsValid, _ := validator.Validate(t.val)
		s.Equalf(t.expValid, actualIsValid, "test case %q", t.name)
	}
}

func (s *ValidatorsSuite) TestValidate_CronValidator() {
	tests := []struct {
		name     string
		val      interface{}
		expValid bool
	}{
		{
			name:     "ok",
			val:      "* * * * *",
			expValid: true,
		},
		{
			name:     "ok macro",
			val:      "@hourly",
			expValid: true,
		},
		{
			name:     "ok empty",
			val:      "",
			expValid: true,
		},
		{
			name:     "misspelled macro",
			val:      "hourly",
			expValid: false,
		},
		{
			name:     "standard missing one element",
			val:      "* * * *",
			expValid: false,
		},
		{
			name:     "one extra element",
			val:      "* * * * * *",
			expValid: false,
		},
		{
			name:     "not string",
			val:      5,
			expValid: false,
		},
	}

	for _, t := range tests {
		valid, err := CronValidator{}.Validate(t.val)
		if t.expValid {
			s.NoError(err)
			s.True(valid)
		} else {
			s.Error(err)
			s.False(valid)
		}
	}
}

func TestValidatorsSuite(t *testing.T) {
	suite.Run(t, new(ValidatorsSuite))
}
