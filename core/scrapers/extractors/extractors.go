package extractors

import (
	"github.com/PuerkitoBio/goquery"
)

type Extractor interface {
	Extract(root *goquery.Selection) (ExtractorResult, error)
}

type ExtractorResult interface{}
