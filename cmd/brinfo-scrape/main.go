package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

var (
	cfgCache             bool
	mergeWithFlag        string
	sourceGUIDFlag       string
	customExtractorsFlag string
	extraDataFlag        string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "brinfo-scraper",
	Short: "CLI for scraping content published by government institutions from Brazil",
}

func main() {
	rootCmd.PersistentFlags().BoolVarP(&cfgCache, "use-cache", "", false, "enable caching, data is kept on .brinfo-cache/")

	rootCmd.AddCommand(scrapeArticleCmd)
	rootCmd.AddCommand(scrapeArticlesListingCmd)

	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
