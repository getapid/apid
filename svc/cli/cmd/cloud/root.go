package cloud

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	apiKeyEnvKey      = "APID_KEY"
	apidCloudEndpoint = "https://api.getapid.com"
)

var (
	apiKey = ""
)

var RootCommand = &cobra.Command{
	Use:   "cloud",
	Short: "Commands for APId cloud",
	Long:  `APId cloud provides a simple and convenient way to monitor your API performance and uptime.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// noop
	},
}

func init() {
	RootCommand.PersistentFlags().StringVarP(&apiKey, "key", "k", os.Getenv(apiKeyEnvKey), "apid access key")
	RootCommand.MarkPersistentFlagRequired("key")
}

func handleError(message string, err error) error {
	fmt.Printf("%s\n\n %s, please try again or report this over at faq.getapid.com\n", err, message)
	return err
}
