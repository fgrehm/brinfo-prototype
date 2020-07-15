package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/fgrehm/brinfo/core"
	op "github.com/fgrehm/brinfo/core/operations"
	xt "github.com/fgrehm/brinfo/core/scrapers/extractors"

	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var scrapeArticleCmd = &cobra.Command{
	Use:   "article [URL]",
	Short: "Scrape contents of articles",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := url.ParseRequestURI(args[0]); err != nil {
			return err
		}

		return runArticleScraper(cmd.Context(), args[0])
	},
}

func init() {
	scrapeArticleCmd.Flags().StringVarP(&mergeWithFlag, "merge-with", "m", "", "JSON to merge with the scraped article data")
	scrapeArticleCmd.Flags().StringVarP(&extraDataFlag, "extra-data", "e", "", "Extra JSON to merge with the scraped article data")
	scrapeArticleCmd.Flags().StringVarP(&sourceGUIDFlag, "source-guid", "s", "", "A string that represents the JSON to merge with the scraped article data")
	scrapeArticleCmd.Flags().StringVarP(&customExtractorsFlag, "custom-extractors", "", "", "A string that represents the JSON of custom extractors to use")
	scrapeArticleCmd.MarkFlagRequired("source-guid")
}

type ArticleData struct {
	*core.ArticleData

	Key    string                 `json:"key"`
	Source string                 `json:"source_guid"`
	Extra  map[string]interface{} `json:"extra,omitempty"`
}

func runArticleScraper(ctx context.Context, url string) error {
	var (
		dataToMerge *core.ArticleData
		extraData   map[string]interface{}
		err         error
		logger      = log.FromContext(ctx)
	)

	if mergeWithFlag != "" {
		dataToMerge, err = core.ArticleDataFromJSON([]byte(mergeWithFlag))
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	if extraDataFlag != "" {
		err = json.Unmarshal([]byte(extraDataFlag), &extraData)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	extractors := []xt.Extractor{xt.BasicArticle()}
	if customExtractorsFlag != "" {
		customExtractors, err := xt.FromJSON([]byte(customExtractorsFlag))
		if err != nil {
			logger.Fatal(err.Error())
		}
		extractors = append(extractors, customExtractors...)
	}

	logger.Infof("Scraping %s", url)
	data, err := op.ScrapeArticle(ctx, op.ScrapeArticleArgs{
		UseCache:   cfgCache,
		URL:        url,
		Extractors: extractors,
		MergeWith:  dataToMerge,
	})
	if err != nil {
		logger.Fatal(err.Error())
	}

	payload := &ArticleData{
		ArticleData: data,
		Extra:       extraData,
		Key:         fmt.Sprintf("%s/article-%s-%s.json", sourceGUIDFlag, data.URLHash, data.FullTextHash),
		Source:      sourceGUIDFlag,
	}
	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		logger.Fatal(err.Error())
	}
	fmt.Println(string(jsonData))

	if valid, msgs := data.ValidForIngestion(); !valid {
		logger.Fatalf("Data is invalid for ingestion: %v", msgs)
	}
	return nil
}
