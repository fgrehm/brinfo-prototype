package operations

import (
	"context"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers"
	. "github.com/fgrehm/brinfo/core/scrapers/extractors"
)

type ScrapeArticleArgs struct {
	UseCache   bool
	URL        string
	Extractors []Extractor
	MergeWith  *ArticleData
}

func ScrapeArticle(ctx context.Context, args ScrapeArticleArgs) (*ArticleData, error) {
	html, httpContentType, err := makeRequest(args.UseCache, args.URL)
	if err != nil {
		return nil, err
	}

	scraper := NewArticleScraper(&ArticleScraperConfig{
		Clock:      &realClock{},
		Extractors: args.Extractors,
		MergeWith:  args.MergeWith,
	})
	return scraper.Run(ctx, html, args.URL, httpContentType)
}
