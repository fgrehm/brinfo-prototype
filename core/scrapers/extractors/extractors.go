package extractors

import (
	"context"

	"github.com/PuerkitoBio/goquery"
)

type ExtractorArgs struct {
	Context         context.Context
	URL             string
	Root            *goquery.Selection
	HTTPContentType string
}

func (a ExtractorArgs) WithRoot(root *goquery.Selection) ExtractorArgs {
	a.Root = root
	return a
}

type Extractor interface {
	Extract(args ExtractorArgs) (ExtractorResult, error)
}

type ExtractorResult interface{}
