package shell

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/getapid/apid/common/variables"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type CommandExecutor interface {
	Exec(string, variables.Variables) ([]byte, error)
}

type ShellCommandExecutor struct {}

func NewShellCommandExecutor() CommandExecutor {
	return &ShellCommandExecutor{}
}

func (e *ShellCommandExecutor) Exec(command string, vars variables.Variables) ([]byte, error) {
	if len(command) == 0 {
		return []byte{}, errors.New("empty command")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sh := os.Getenv("SHELL")

	cmd := exec.CommandContext(ctx, sh, "-c", command)
	cmd.Env = append(os.Environ(), getEnvFromVars(vars)...)
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	res := out.String()
	return 	[]byte(strings.ReplaceAll(res, "\n", "")), err
}

func getEnvFromVars(vars variables.Variables) []string {
	var result []string
	for key, value := range vars.Raw() {
		result = append(result, flattenVars(strings.ToUpper(key), value)...)
	}
	return result
}

func flattenVars(namespace string, vars interface{}) []string {
	switch val := vars.(type) {
	case map[string]interface{}:
		var result []string
		for key, value := range val {
			result = append(result, flattenVars(strings.ToUpper(namespace + "_" + key), value)...)
		}
		return result
	case []interface{}:
		var result []string
		for index, value := range val {
			result = append(result, fmt.Sprintf("%s=%v", strings.ToUpper(namespace + "_" + strconv.Itoa(index)), value))
		}
		return result
	default:
		return []string{fmt.Sprintf("%s=%v", namespace, val)}
	}
}