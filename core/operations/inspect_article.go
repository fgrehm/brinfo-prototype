package operations

import (
	"context"
	"errors"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers"
)

type InspectArticleArgs struct {
	UseCache          bool
	ContentSourceRepo ContentSourceRepo
	URL               string
	ArticleScraper    ArticleScraper
	MergeWith         string
}

type InspectedArticleData struct {
	*ScrapedArticleData
}

func InspectArticle(ctx context.Context, args InspectArticleArgs) (*InspectedArticleData, error) {
	var (
		url       = args.URL
		scraper   = args.ArticleScraper
		cache     = args.UseCache
		mergeWith *ScrapedArticleData

		cs  *ContentSource
		err error
		log = loggerFromContext(ctx)
	)

	if args.MergeWith != "" {
		mergeWith, err = ScrapedArticleDataFromJSON([]byte(args.MergeWith))
		log.Debugf("Will merge with '%s'", args.MergeWith)
		if err != nil {
			return nil, err
		}
	}
	if err = validateInspectArticleArgs(args); err != nil {
		return nil, err
	}

	if args.ContentSourceRepo != nil {
		log.Debugf("Looking up content source for %s", args.URL)
		cs, err = lookupContentSourceForUrl(ctx, args.ContentSourceRepo, args.URL)
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

	data, err := doScrapeArticle(ctx, cache, url, cs, scraper)
	if err != nil {
		return nil, err
	}

	if mergeWith != nil {
		log.Debugf("Merging with %s", args.MergeWith)
		data.CollectValues(mergeWith)
	}

	return &InspectedArticleData{data}, nil
}

func validateInspectArticleArgs(args InspectArticleArgs) error {
	if args.URL == "" {
		return errors.New("No URL provided")
	}
	if args.ContentSourceRepo == nil {
		return errors.New("No content source repo provided")
	}

	return nil
}
