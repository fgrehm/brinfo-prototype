package extractors

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type structuredExtractor struct {
	selector string
	configs  map[string]Extractor
	multiple bool
}

func Structured(selector string, extractors map[string]Extractor) Extractor {
	if selector == "" {
		panic("No selector provided")
	}
	if len(extractors) == 0 {
		panic("No extractors provided")
	}
	return &structuredExtractor{selector, extractors, false}
}

func StructuredList(selector string, extractors map[string]Extractor) Extractor {
	if selector == "" {
		panic("No selector provided")
	}
	if len(extractors) == 0 {
		panic("No extractors provided")
	}
	return &structuredExtractor{selector, extractors, true}
}

func (e *structuredExtractor) Extract(args ExtractorArgs) (ExtractorResult, error) {
	root := args.Root.Find(e.selector)

	if root.Length() == 0 {
		return nil, fmt.Errorf("'%s' not found", e.selector)
	}

	if !e.multiple {
		if root.Length() > 1 {
			return nil, fmt.Errorf("Multiple '%s' found (%d)", e.selector, args.Root.Length())
		}

		if root.Length() == 1 {
			return e.extractOne(args.WithRoot(root))
		}
	}

	var err error
	result := []map[string]ExtractorResult{}

	root.Each(func(idx int, s *goquery.Selection) {
		if err != nil {
			return
		}

		value, innerErr := e.extractOne(args.WithRoot(s))
		if innerErr != nil {
			err = innerErr
			return
		}

		result = append(result, value)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e *structuredExtractor) extractOne(args ExtractorArgs) (map[string]ExtractorResult, error) {
	ret := map[string]ExtractorResult{}
	for fieldName, extractor := range e.configs {
		result, err := extractor.Extract(args)
		if err != nil {
			// errors.WithStack
			return nil, fmt.Errorf("within '%s' > %s", e.selector, err)
		}
		ret[fieldName] = result
	}

	return ret, nil
}
