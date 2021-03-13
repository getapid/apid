package cloud

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var enableSuite = &cobra.Command{
	Use:     "enable",
	Short:   "Enables a suite",
	Long:    `Enables a suite for execution with the predefined schedule.`,
	Example: `apid cloud suite enable -n <suite name> -k <api key>`,
	Args:    cobra.NoArgs,
	Run:     enableSuiteCommandFunc,
}

func init() {
	suiteCommand.AddCommand(enableSuite)

	enableSuite.Flags().StringVarP(&suiteName, "name", "n", "", "suite name, required")
	enableSuite.MarkFlagRequired("name")
}

func setSuiteStatus(suiteID string, enabled bool) error {
	status := "disable"
	if enabled {
		status = "enable"
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/cli/suite/%s/%s", apidCloudEndpoint, suiteID, status), nil)
	if err != nil {
		return handleError("unexcpected error encountered while preparing apid cloud request", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return handleError("could not connect to apid cloud", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleError(fmt.Sprintf("could not %s suite to apid cloud, api returned an error", status), errors.New("apid_cloud_internal_server_error"))
	}
	fmt.Println("success")
	return nil
}

func enableSuiteCommandFunc(cmd *cobra.Command, args []string) {
	suites, err := getSuites()
	if err != nil {
		return
	}
	for _, suite := range suites {
		if suite.Name == suiteName {
			setSuiteStatus(suite.SuiteID, true)
			return
		}
	}

	fmt.Println("suite not found")
}
