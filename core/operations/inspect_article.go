package operations

import (
	"errors"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers"
)

type InspectArticleInput struct {
	ContentSourceRepo ContentSourceRepo
	Url               string
	ArticleScraper    ArticleScraper
}

type InspectedArticleData struct {
	*ScrapedArticleData
}

func InspectArticle(input InspectArticleInput) (*InspectedArticleData, error) {
	var (
		url     = input.Url
		scraper = input.ArticleScraper

		cs  *ContentSource
		err error
	)

	if err = validateInspectArticleInput(input); err != nil {
		return nil, err
	}

	if input.ContentSourceRepo != nil {
		cs, err = lookupContentSourceForUrl(input.ContentSourceRepo, input.Url)
		if err != nil {
			return nil, err
		}
	}

	if cs != nil {
		if err = validateContentSourceForScraping(cs, url); err != nil {
			return nil, err
		}
	}

	if scraper == nil && cs != nil {
		scraper = cs.ArticleScraper
	}
	if scraper == nil {
		scraper = DefaultArticleScraper
	}

	data, err := doScrapeArticleContent(url, cs, scraper)
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
