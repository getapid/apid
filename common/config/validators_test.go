package config

import (
	"testing"

	"github.com/getapid/apid-cli/common/step"
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
			name: "valid",
			val: &step.ExpectBody{
				Type:    pstring("json"),
				Exact:   pbool(true),
				Content: "",
			},
			expValid: true,
		},
		{
			name:     "defaults",
			val:      &step.ExpectBody{},
			expValid: true,
		},
		{
			name: "valid - plaintext and non exact",
			val: &step.ExpectBody{
				Type:  pstring("plaintext"),
				Exact: pbool(true),
			},
			expValid: true,
		},
		{
			name: "valid - plaintext and non exact",
			val: &step.ExpectBody{
				Exact: pbool(false),
			},
			expValid: true,
		},
		{
			name: "invalid type",
			val: &step.ExpectBody{
				Type: pstring("wqieufwiu"),
			},
			expValid: false,
		},
	}

	validator := ExpectBodyValidator{}
	for _, t := range testCases {
		actualIsValid, actualErr := validator.Validate(t.val)

		s.Equalf(t.expValid, actualIsValid, "test case %q", t.name)
		if t.expValid {
			s.NoErrorf(actualErr, "test case %q", t.name)
			if t.val != nil { // check default were set
				s.NotNilf(t.val.(*step.ExpectBody).Type, "test case %q", t.name)
				s.NotNilf(t.val.(*step.ExpectBody).Exact, "test case %q", t.name)
			}
		} else {
			s.Errorf(actualErr, "test case %q", t.name)
		}
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
