package transaction

import (
	"github.com/iv-p/apiping/svc/client/step"
	"github.com/iv-p/apiping/svc/client/variables"
)

// Transaction is the definition of a transaction
type Transaction struct {
	ID        string              `yaml:"id"`
	Variables variables.Variables `yaml:"variables"`
	Steps     []step.Step         `yaml:"steps"`
}
