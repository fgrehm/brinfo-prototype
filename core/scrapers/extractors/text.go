package extractors

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type textExtractor struct {
	selector string
	multiple bool
	required bool
}

func Text(selector string, multiple bool) Extractor {
	return &textExtractor{
		selector: selector,
		multiple: multiple,
		required: true,
	}
}

func OptText(selector string, multiple bool) Extractor {
	return &textExtractor{
		selector: selector,
		multiple: multiple,
		required: false,
	}
}

func (e *textExtractor) Extract(root *goquery.Selection) (ExtractorResult, error) {
	sel := root.Find(e.selector)

	if sel.Length() == 0 {
		if e.required {
			return nil, fmt.Errorf("'%s' not found", e.selector)
		} else {
			return nil, nil
		}
	}

	if !e.multiple {
		if sel.Length() > 1 {
			return nil, fmt.Errorf("Multiple '%s' found (%d)", e.selector, sel.Length())
		}

		if sel.Length() == 1 {
			text := sel.Text()
			if text == "" {
				if e.required {
					return nil, fmt.Errorf("'%s' is empty", e.selector)
				} else {
					return nil, nil
				}
			}
			return strings.Trim(text, " "), nil
		}
	}

	var err error
	ret := sel.Map(func(idx int, s *goquery.Selection) string {
		return strings.Trim(s.Text(), " ")
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}
