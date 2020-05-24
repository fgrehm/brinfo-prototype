package operations

import (
	"errors"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers"
)

type InspectBytesInput struct {
	Html           []byte
	ArticleScraper ArticleScraper
	Url            string
}

func InspectBytes(input InspectBytesInput) (interface{}, error) {
	if len(input.Html) == 0 {
		return nil, errors.New("No HTML provided")
	}
	if len(input.Url) == 0 {
		return nil, errors.New("No Url provided")
	}

	scraper := input.ArticleScraper
	if scraper == nil {
		// scraper = CombinedScraper(DefaultArticleScraper, ...)
		scraper = DefaultArticleScraper
	}

	data, err := scraper.Run(input.Html, input.Url, `text/html; charset="UTF-8"`)
	if err != nil {
		return nil, err
	}

	return data, nil
}
