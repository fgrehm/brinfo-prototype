package operations

import (
	"errors"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers"
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

	body, contentType, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	scraper := input.ArticleScraper
	if scraper == nil {
		// scraper = CombinedScraper(DefaultArticleScraper, ...)
		scraper = DefaultArticleScraper
	}

	data, err := scraper.Run(body, url, contentType)
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
