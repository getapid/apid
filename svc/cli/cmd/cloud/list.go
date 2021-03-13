package cloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List suites",
	Long:    `Get all suites associated with this api key.`,
	Example: `apid cloud suite list -k <api-key>`,
	Args:    cobra.NoArgs,
	Run:     listCommandFunc,
}

type getSuitesListResponse struct {
	Suites []suite `json:"suites"`
}

type suite struct {
	SuiteID string `json:"suite_id"`
	Name    string `json:"name"`
}

func init() {
	suiteCommand.AddCommand(listCommand)
}

func getSuites() ([]suite, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cli/suite", apidCloudEndpoint), nil)
	if err != nil {
		handleError("unexcpected error encountered while preparing apid cloud request", err)
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		handleError("could not connect to apid cloud", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError("unexcpected error encountered while receiving response from apid cloud", err)
		return nil, err
	}
	var data getSuitesListResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		handleError("unexcpected error encountered while transforming response from apid cloud", err)
	}
	return data.Suites, err
}

func listCommandFunc(cmd *cobra.Command, args []string) {
	if len(apiKey) == 0 {
		log.Fatal("missing api key")
	}

	suites, err := getSuites()
	if err != nil {
		return
	}

	for _, suite := range suites {
		fmt.Println(suite.Name)
	}
}
