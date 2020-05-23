package operations

import (
	"bytes"
	"errors"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers"

	"github.com/dimchansky/utfbom"
)

type InspectArticleInput struct {
	Url            string
	ArticleScraper ArticleScraper
}

type InspectedArticleData struct {
	*ScrapedArticleData
}

func InspectArticle(input InspectArticleInput) (*InspectedArticleData, error) {
	var (
		url = input.Url

		err error
	)

	if err = validateInspectArticleInput(input); err != nil {
		return nil, err
	}

	body, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	scraper := input.ArticleScraper
	if scraper == nil {
		// scraper = CombinedScraper(DefaultArticleScraper, ...)
		scraper = DefaultArticleScraper
	}

	buf := bytes.NewBuffer(body)
	data, err := scraper.Run(utfbom.SkipOnly(buf), url)
	if err != nil {
		return nil, err
	}

	return &InspectedArticleData{data}, nil
}

func validateInspectArticleInput(input InspectArticleInput) error {
	if input.Url == "" {
		return errors.New("No URL provided")
	}

	return nil
}
