package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/fgrehm/brinfo/core"
	op "github.com/fgrehm/brinfo/core/operations"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var scrapeCmd = &cobra.Command{
	Use:   "scrape-article [URL]",
	Short: "Scrape contents of articles from well known sources",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		urlToScrape, err := url.ParseRequestURI(args[0])
		if err != nil {
			return err
		}

		log.Infof("Scraping %s", urlToScrape)
		data, err := op.ScrapeArticleContent(op.ScrapeArticleContentInput{
			Url:  urlToScrape.String(),
			Repo: repo,
		})
		if err != nil {
			log.Fatal(err)
		}
		if !data.ValidForIngestion() {
			log.Fatal("Data is invalid for ingestion")
		}

		payload := &ArticleDataToUpload{data, fmt.Sprintf("%s/article-%s.json", data.SourceID, data.UrlHash)}
		jsonData, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		outputDir := fmt.Sprintf("./output/%s", data.SourceID)
		os.MkdirAll(outputDir, os.ModePerm)

		outputPath := fmt.Sprintf("%s/article-%s.json", outputDir, data.UrlHash)
		log.Infof("Saving to %s", outputPath)
		if err = ioutil.WriteFile(outputPath, jsonData, 0644); err != nil {
			log.Fatal(err)
		}

		// fmt.Println(string(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
}

type ArticleDataToUpload struct {
	*core.ScrapedArticleData

	Key string `json:"key"`
}
