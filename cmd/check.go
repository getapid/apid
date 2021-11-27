package cmd

import (
	"github.com/getapid/apid/file"
	"github.com/getapid/apid/http"
	"github.com/getapid/apid/spec"
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

func check(cmd *cobra.Command, args []string) error {
	specLoader := spec.NewSpecLoader(file.JsonnetReader{})
	specs := specLoader.Load(specPattern)

	var w writer.Writer
	if json {
		w = writer.NewJSON(cmd.OutOrStdout())
	} else {
		w = writer.NewConsole(cmd.OutOrStdout(), silent)
	}

	httpClient := http.NewClient()
	stepHttpClient := step.NewHTTPClient(httpClient)
	stepInterpolator := step.NewInterpolator()
	stepValidator := step.NewValidator()
	stepExporter := step.NewExporter()
	stepRunner := step.NewRunner(stepHttpClient, *stepInterpolator, stepValidator, stepExporter)
	specRunner := runner.NewParallelSpecRunner(parallelism, stepRunner, w)

	w.Prelude()

	specRunner.Run(specs)

	w.Conclusion()
	return nil
}
