package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	op "github.com/fgrehm/brinfo/core/operations"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var inspectArticle = &cobra.Command{
	Use:   "inspect-article [URL]",
	Short: "Inspect an article's page and output some useful information about it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		urlToInspect, err := url.ParseRequestURI(args[0])
		if err != nil {
			return err
		}

		logger := log.FromContext(cmd.Context())
		logger.Infof("Inspecting %s", urlToInspect)
		data, err := op.InspectArticle(cmd.Context(), op.InspectArticleInput{
			Url:               urlToInspect.String(),
			ContentSourceRepo: repo,
		})
		if err != nil {
			panic(err)
		}

		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))

		if !data.ValidForIngestion() {
			logger.Error("Won't be able to ingest article")
			panic("aborting")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(inspectArticle)
}
