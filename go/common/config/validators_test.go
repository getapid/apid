package config

import (
	"testing"

	"github.com/iv-p/apid/common/step"

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
			name: "invalid - plaintext and non exact",
			val: &step.ExpectBody{
				Type:  pstring("plaintext"),
				Exact: pbool(false),
			},
			expValid: false,
		},
		{
			name: "invalid - plaintext and non exact",
			val: &step.ExpectBody{
				Exact: pbool(false),
			},
			expValid: false,
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

func TestValidatorsSuite(t *testing.T) {
	suite.Run(t, new(ValidatorsSuite))
}
