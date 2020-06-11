package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	op "github.com/fgrehm/brinfo/core/operations"

	"github.com/spf13/cobra"
)

var scrapeArticlesListingFlags = struct {
	linkContainer        string
	urlExtractor         string
	publishedAtExtractor string
	imageURLExtractor    string
}{}

var scrapeArticlesListingCmd = &cobra.Command{
	Use:   "scrape-articles-listing [URL]",
	Short: "Extract a list of article links and metadata from a page that has a list of articles",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := url.ParseRequestURI(args[0])
		if err != nil {
			return err
		}

		data, err := op.ScrapeArticlesListing(cmd.Context(), op.ScrapeArticlesListingArgs{
			URL:                  args[0],
			LinkContainer:        scrapeArticlesListingFlags.linkContainer,
			URLExtractor:         scrapeArticlesListingFlags.urlExtractor,
			PublishedAtExtractor: scrapeArticlesListingFlags.publishedAtExtractor,
			ImageURLExtractor:    scrapeArticlesListingFlags.imageURLExtractor,
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
	rootCmd.AddCommand(scrapeArticlesListingCmd)
	scrapeArticlesListingCmd.Flags().StringVarP(&scrapeArticlesListingFlags.linkContainer, "link-container", "l", "", "CSS selector for the element that wraps links to scrape (required)")
	scrapeArticlesListingCmd.MarkFlagRequired("link-container")
	scrapeArticlesListingCmd.Flags().StringVarP(&scrapeArticlesListingFlags.urlExtractor, "url-extractor", "u", "a[href] | href", "CSS selector for the actual link, nested under the elements wrapped by the container")
	scrapeArticlesListingCmd.Flags().StringVarP(&scrapeArticlesListingFlags.publishedAtExtractor, "published-at-extractor", "p", "", "CSS selector for the actual link, nested under the elements wrapped by the container")
	scrapeArticlesListingCmd.Flags().StringVarP(&scrapeArticlesListingFlags.imageURLExtractor, "image-url-extractor", "i", "", "CSS selector for the actual link, nested under the elements wrapped by the container")
}
