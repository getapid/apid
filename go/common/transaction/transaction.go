package transaction

import (
	"github.com/iv-p/apid/common/step"
)

// Transaction is the definition of a transaction
type Transaction struct {
	ID        string                 `yaml:"id"`
	Variables map[string]interface{} `yaml:"variables"`
	Steps     []step.Step            `yaml:"steps"`
}

// Result holds information for the result of a transaction
// Passed to the result writer. Content TBD
type Result struct{}
