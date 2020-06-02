package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	op "github.com/fgrehm/brinfo/core/operations"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var inspectFile = &cobra.Command{
	Use:   "inspect-file [FILE] [URL]",
	Short: "Inspect an HTML page saved locally using the default scraper",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		logger := log.FromContext(cmd.Context())

		logger.Infof("Reading %s", args[0])
		fileContents, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}

		data, err := op.InspectBytes(cmd.Context(), op.InspectBytesInput{
			Html: fileContents,
			Url:  args[1],
		})
		if err != nil {
			panic(err)
		}

		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectFile)
}
