package extractors

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type structuredExtractor struct {
	selector string
	configs  map[string]Extractor
}

func Structured(selector string, extractors map[string]Extractor) Extractor {
	if selector == "" {
		panic("No selector provided")
	}
	if len(extractors) == 0 {
		panic("No extractors provided")
	}
	return &structuredExtractor{selector, extractors}
}

func (e *structuredExtractor) Extract(sel *goquery.Selection) (ExtractorResult, error) {
	ret := map[string]ExtractorResult{}
	root := sel.Find(e.selector)

	// TODO: Error if multiple found

	// TODO: Collect all errors instead of returning on the first one
	for fieldName, extractor := range e.configs {
		result, err := extractor.Extract(root)
		if err != nil {
			return nil, fmt.Errorf("within '%s' > %s", e.selector, err)
		}
		ret[fieldName] = result
	}

	return ret, nil
}
