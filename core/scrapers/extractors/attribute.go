package extractors

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type attrExtractor struct {
	selector string
	attr     string
	multiple bool
}

func Attribute(selector, attr string) Extractor {
	return &attrExtractor{selector, attr, false}
}

func (e *attrExtractor) Extract(root *goquery.Selection) (ExtractorResult, error) {
	sel := root.Find(e.selector)
	if sel.Length() == 0 {
		return nil, fmt.Errorf("'%s' not found", e.selector)
	}

	if !e.multiple {
		if sel.Length() > 1 {
			return nil, fmt.Errorf("Multiple '%s' found (%d)", e.selector, sel.Length())
		}

		if sel.Length() == 1 {
			attr, found := sel.Attr(e.attr)
			if !found {
				return nil, fmt.Errorf("Attribute '%s' for '%s' not found", e.attr, e.selector)
			}
			return attr, nil
		}
	}
	if sel.Length() > 1 && !e.multiple {
		return nil, fmt.Errorf("Multiple '%s' found (%d)", e.selector, sel.Length())
	}

	var err error
	ret := sel.Map(func(idx int, s *goquery.Selection) string {
		attr, found := s.Attr(e.attr)
		if !found {
			err = fmt.Errorf("Attribute '%s' for '%s'[%d] not found", e.attr, e.selector, idx)
			return ""
		}
		return attr
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}
