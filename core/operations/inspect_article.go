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
	MergeWith         string
}

type InspectedArticleData struct {
	*ScrapedArticleData
}

func InspectArticle(ctx context.Context, input InspectArticleInput) (*InspectedArticleData, error) {
	var (
		url       = input.Url
		scraper   = input.ArticleScraper
		mergeWith *ScrapedArticleData

		cs  *ContentSource
		err error
		log = loggerFromContext(ctx)
	)

	if input.MergeWith != "" {
		mergeWith, err = ScrapedArticleDataFromJSON([]byte(input.MergeWith))
		log.Debugf("Will merge with '%s'", input.MergeWith)
		if err != nil {
			return nil, err
		}
	}
	if err = validateInspectArticleInput(input); err != nil {
		return nil, err
	}

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

	if mergeWith != nil {
		log.Debugf("Merging with %s", input.MergeWith)
		data.CollectValues(mergeWith)
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
