package operations

import (
	. "github.com/fgrehm/brinfo/core"
)

type ScrapeArticleContentInput struct {
	Repo          ContentSourceRepo
	ContentSource *ContentSource
	Url           string
}

func ScrapeArticleContent(input ScrapeArticleContentInput) (*ScrapedArticleData, error) {
	var (
		url  = input.Url
		cs   = input.ContentSource
		repo = input.Repo

		err error
	)

	if err = validateScrapeArticleContentInput(input); err != nil {
		return nil, err
	}

	if cs == nil {
		cs, err = mustLookupContentSourceForUrl(repo, url)
		if err != nil {
			return nil, err
		}
	}

	if err = validateContentSourceForScraping(cs, url); err != nil {
		return nil, err
	}

	data, err := doScrapeArticleContent(url, cs, cs.ArticleScraper)
	if err != nil {
		return nil, err
	}
	return data, nil
}
