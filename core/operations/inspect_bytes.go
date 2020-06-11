package operations

import (
	"context"
	"errors"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers"
)

type InspectBytesArgs struct {
	HTML           []byte
	URL            string
	ArticleScraper ArticleScraper
	ContentType    string
}

func InspectBytes(ctx context.Context, args InspectBytesArgs) (*ScrapedArticleData, error) {
	if len(args.HTML) == 0 {
		return nil, errors.New("No HTML provided")
	}
	if len(args.URL) == 0 {
		return nil, errors.New("No URL provided")
	}

	if args.ContentType == "" {
		args.ContentType = "text/html; charset=UTF-8"
	}
	if args.ArticleScraper == nil {
		args.ArticleScraper = DefaultArticleScraper
	}

	data, err := args.ArticleScraper.Run(ctx, args.HTML, args.URL, args.ContentType)
	if err != nil {
		return nil, err
	}

	return data, nil
}
