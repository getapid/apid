package transaction

import (
	"github.com/iv-p/apid/common/step"
)

// Transaction is the definition of a transaction
type Transaction struct {
	ID        string                 `yaml:"id" validate:"required"`
	Variables map[string]interface{} `yaml:"variables"`
	Steps     []step.Step            `yaml:"steps" validate:"required,unique=ID"`
}

// TODO remove after rebase
// Result holds information for the result of a transaction
// Passed to the result writer. Content TBD
type Result struct{}
