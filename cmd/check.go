package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/getapid/apid/env"
	"github.com/getapid/apid/http"
	"github.com/getapid/apid/log"
	"github.com/getapid/apid/spec"
	"github.com/getapid/apid/spec/loader"
	"github.com/getapid/apid/spec/runner"
	"github.com/getapid/apid/step"
	"github.com/getapid/apid/writer"
	"github.com/spf13/cobra"
)

var (
	specPattern string
	parallelism int
	json        bool
	silent      bool

	ErrTestFailedError = errors.New("tests_failed")
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Runs one or more specs",
	Long: `Loads a spec file. In the case a file glob is provided
	 (see examples below) checks all specs matching that pattern.`,
	Example: `
	apid check	
	apid check --specs spec.jsonnet
	apid check -s tests/**/*.jsonnet`,
	Args: cobra.NoArgs,
	RunE: check,
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&specPattern, "specs", "s", "spec.jsonnet", "specification to run")
	checkCmd.Flags().IntVarP(&parallelism, "parallelism", "p", 10, "number of concurrent transaction executions")
	checkCmd.Flags().BoolVar(&json, "json", false, "set output mode to json")
	checkCmd.Flags().BoolVar(&silent, "silent", false, "set output mode to silent only printing the end result of the tests")
}

const (
	ExitCodeOK int = 0

	ExitCodeErrTestFailure       int = 1
	ExitCodeErrReadingSpec       int = 100
	ExitCodeErrInvalidSpecFormat int = 101
)

func check(cmd *cobra.Command, args []string) error {
	files, err := filepath.Glob(specPattern)
	if err != nil {
		log.L.Errorf("error finding files: %s", err)
	}

	if len(files) == 0 {
		log.L.Fatalf("no files found matching pattern: %s", specPattern)
	}

	var specLoader loader.JsonnetLoader
	hasError := false

	specs := make(map[string]spec.Spec)
	for _, file := range files {
		jsonSpecs, err := specLoader.Load(file, env.LoadVars())
		if err != nil {
			os.Exit(ExitCodeErrReadingSpec)
		}
		for name, jsonSpec := range jsonSpecs {
			spec, err := spec.Unmarshal([]byte(jsonSpec))
			hasError = hasError || err != nil

			if _, ok := specs[name]; ok {
				log.L.Fatalf("duplicate spec with name %s", name)
			}

			if len(spec.Name) == 0 {
				spec.Name = name
			}

			specs[name] = spec
		}
	}

	if hasError {
		os.Exit(ExitCodeErrInvalidSpecFormat)
	}

	var w writer.Writer
	if json {
		w = writer.NewJSON(cmd.OutOrStdout())
	} else {
		w = writer.NewConsole(cmd.OutOrStdout(), silent)
	}
	w.Prelude()

	httpClient := http.NewClient()
	stepHttpClient := step.NewHTTPClient(httpClient)
	stepInterpolator := step.NewInterpolator()
	stepValidator := step.NewValidator()
	stepExporter := step.NewExporter()
	stepRunner := step.NewRunner(stepHttpClient, *stepInterpolator, stepValidator, stepExporter)
	specRunner := runner.NewParallelSpecRunner(parallelism, stepRunner, w)
	specRunner.Run(specs)

	w.Conclusion()
	return nil
}
