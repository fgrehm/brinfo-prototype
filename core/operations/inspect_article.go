package operations

import (
	"context"
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

func InspectArticle(ctx context.Context, input InspectArticleInput) (*InspectedArticleData, error) {
	var (
		url     = input.Url
		scraper = input.ArticleScraper

		cs  *ContentSource
		err error
	)

	if err = validateInspectArticleInput(input); err != nil {
		return nil, err
	}

	log := loggerFromContext(ctx)
	if input.ContentSourceRepo != nil {
		log.Debugf("Looking up content source for %s", input.Url)
		cs, err = lookupContentSourceForUrl(ctx, input.ContentSourceRepo, input.Url)
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
		log.Debugf("Using custom scraper for %s", cs.ID)
		scraper = cs.ArticleScraper
	}
	if scraper == nil {
		log.Debug("Using default article scraper")
		scraper = DefaultArticleScraper
	}

	data, err := doScrapeArticle(ctx, url, cs, scraper)
	if err != nil {
		return nil, err
	}

	return &InspectedArticleData{data}, nil
}

func validateInspectArticleInput(input InspectArticleInput) error {
	if input.Url == "" {
		return errors.New("No URL provided")
	}
	if input.ContentSourceRepo == nil {
		return errors.New("No content source repo provided")
	}

	return nil
}
