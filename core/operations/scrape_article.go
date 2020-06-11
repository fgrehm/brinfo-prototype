package operations

import (
	"context"
	"errors"

	. "github.com/fgrehm/brinfo/core"
)

type ScrapeArticleInput struct {
	UseCache      bool
	Repo          ContentSourceRepo
	ContentSource *ContentSource
	Url           string
}

func ScrapeArticle(ctx context.Context, input ScrapeArticleInput) (*ScrapedArticleData, error) {
	var (
		url   = input.Url
		cs    = input.ContentSource
		repo  = input.Repo
		cache = input.UseCache

		err error
	)

	if err = input.validate(); err != nil {
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

func (i ScrapeArticleInput) validate() error {
	if i.Url == "" {
		return errors.New("No URL provided")
	}
	if i.ContentSource == nil && i.Repo == nil {
		return errors.New("No ContentSource or Repository provided")
	}

	return nil
}
