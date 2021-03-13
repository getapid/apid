package cloud

import (
	"fmt"

	"github.com/spf13/cobra"
)

var disableSuite = &cobra.Command{
	Use:     "disable",
	Short:   "Disables a suite",
	Long:    `Disables a suite for execution. Does not stop current execution if one is in progress`,
	Example: `apid cloud suite disable -n <suite name> -k <api key>`,
	Args:    cobra.NoArgs,
	Run:     disableSuiteCommandFunc,
}

func init() {
	suiteCommand.AddCommand(disableSuite)

	disableSuite.Flags().StringVarP(&suiteName, "name", "n", "", "suite name, required")
	disableSuite.MarkFlagRequired("name")
}

func disableSuiteCommandFunc(cmd *cobra.Command, args []string) {
	suites, err := getSuites()
	if err != nil {
		return
	}
	for _, suite := range suites {
		if suite.Name == suiteName {
			setSuiteStatus(suite.SuiteID, false)
			return
		}
	}

	fmt.Println("suite not found")
}
