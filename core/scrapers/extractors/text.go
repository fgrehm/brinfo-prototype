package extractors

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type textExtractor struct {
	selector string
	multiple bool
}

func Text(selector string, multiple bool) Extractor {
	return &textExtractor{selector, multiple}
}

func (e *textExtractor) Extract(root *goquery.Selection) (ExtractorResult, error) {
	sel := root.Find(e.selector)

	if sel.Length() == 0 {
		return nil, fmt.Errorf("'%s' not found", e.selector)
	}

	if !e.multiple {
		if sel.Length() > 1 {
			return nil, fmt.Errorf("Multiple '%s' found (%d)", e.selector, sel.Length())
		}

		if sel.Length() == 1 {
			return sel.Text(), nil
		}
	}

	var err error
	ret := sel.Map(func(idx int, s *goquery.Selection) string {
		return s.Text()
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}
