package cloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/getapid/apid-cli/common/config"
	"github.com/spf13/cobra"
)

var (
	configFilepath string
	suiteName      = ""
)

type upsertSuiteRequest struct {
	Name   string        `json:"name"`
	Config config.Config `json:"config"`
}

var upsertCommand = &cobra.Command{
	Use:   "upsert",
	Short: "Upserts a suite",
	Long: `Creates a new suite with the corresponding name. If the name already
	exists the suite will be replaced.`,
	Example: `apid cloud suite upsert -c <config directory/file> -n <suite name> -k <api key>`,
	Args:    cobra.NoArgs,
	Run:     upsertCommandFunc,
}

func init() {
	suiteCommand.AddCommand(upsertCommand)

	upsertCommand.Flags().StringVarP(&configFilepath, "config", "c", "./apid.yaml", "file with config to run")
	upsertCommand.MarkFlagRequired("config")

	upsertCommand.Flags().StringVarP(&suiteName, "name", "n", "", "suite name, required")
	upsertCommand.MarkFlagRequired("name")
}

func upsertSuite(cfg config.Config) error {
	serialized, err := json.Marshal(upsertSuiteRequest{
		Name:   suiteName,
		Config: cfg,
	})
	if err != nil {
		return handleError("unexcpected error encountered while serializing suite", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/cli/suite", apidCloudEndpoint), bytes.NewBuffer(serialized))
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
		return handleError("could not upsert suite to apid cloud, api returned an error", errors.New("apid_cloud_internal_server_error"))
	}
	fmt.Println("success")
	return nil
}

func upsertCommandFunc(cmd *cobra.Command, args []string) {
	c, err := config.Load(configFilepath)
	if err != nil {
		handleError("could not load config", err)
		return
	}

	err = config.Validate(c)
	if err != nil {
		handleError("suite failed validation", err)
		return
	}

	upsertSuite(c)
}
