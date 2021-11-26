package env

import (
	"os"
	"strings"
)

type Vars map[string]string

func LoadVars() Vars {
	environ := os.Environ()
	env := make(Vars, len(environ))
	for _, e := range environ {
		pair := strings.SplitN(e, "=", 2)
		env[pair[0]] = pair[1]
	}
	return env
}
