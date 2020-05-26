package operations

import (
	"errors"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers"
)

type InspectBytesInput struct {
	Html              []byte
	ContentSourceRepo ContentSourceRepo
	ContentSource     *ContentSource
	ArticleScraper    ArticleScraper
	Url               string
	ContentType       *string
}

func InspectBytes(input InspectBytesInput) (interface{}, error) {
	if len(input.Html) == 0 {
		return nil, errors.New("No HTML provided")
	}
	if len(input.Url) == 0 {
		return nil, errors.New("No Url provided")
	}

	cs, err := fetchContentSource(input)
	if err != nil {
		return nil, err
	}

	contentType := "text/html; charset=UTF-8"
	if cs != nil && cs.ForceContentType != "" {
		contentType = cs.ForceContentType
	}

	scraper := fetchScraper(input, cs)
	data, err := scraper.Run(input.Html, input.Url, contentType)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func fetchContentSource(input InspectBytesInput) (*ContentSource, error) {
	if input.ContentSource != nil {
		return input.ContentSource, nil
	}

	if input.ContentSourceRepo != nil {
		return lookupContentSourceForUrl(input.ContentSourceRepo, input.Url)
	}

	return nil, nil
}

func fetchScraper(input InspectBytesInput, cs *ContentSource) ArticleScraper {
	if input.ArticleScraper != nil {
		return input.ArticleScraper
	}

	if cs != nil && cs.ArticleScraper != nil {
		return cs.ArticleScraper
	}

	return DefaultArticleScraper
}
