package operations

import (
	"context"
	"errors"

	. "github.com/fgrehm/brinfo/core"
)

type ScrapeArticleArgs struct {
	UseCache      bool
	Repo          ContentSourceRepo
	ContentSource *ContentSource
	URL           string
}

func ScrapeArticle(ctx context.Context, args ScrapeArticleArgs) (*ScrapedArticleData, error) {
	var (
		url   = args.URL
		cs    = args.ContentSource
		repo  = args.Repo
		cache = args.UseCache

		err error
	)

	if err = args.validate(); err != nil {
		return nil, err
	}

	if cs == nil {
		cs, err = mustLookupContentSourceForUrl(ctx, repo, url)
		if err != nil {
			return nil, err
		}
	}

	if err = validateContentSourceForScraping(cs, url); err != nil {
		return nil, err
	}

	data, err := doScrapeArticle(ctx, cache, url, cs, cs.ArticleScraper)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (i ScrapeArticleArgs) validate() error {
	if i.URL == "" {
		return errors.New("No URL provided")
	}
	if i.ContentSource == nil && i.Repo == nil {
		return errors.New("No ContentSource or Repository provided")
	}

	return nil
}
