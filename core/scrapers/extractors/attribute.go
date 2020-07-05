package extractors

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type attrExtractor struct {
	selector string
	attr     string
	multiple bool
	required bool
}

func Attribute(selector, attr string) Extractor {
	return &attrExtractor{
		selector: selector,
		attr:     attr,
		multiple: false,
		required: true,
	}
}

func OptAttribute(selector, attr string) Extractor {
	return &attrExtractor{
		selector: selector,
		attr:     attr,
		multiple: false,
		required: false,
	}
}

func (e *attrExtractor) Extract(args ExtractorArgs) (ExtractorResult, error) {
	sel := args.Root.Find(e.selector)
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
			attr, found := sel.Attr(e.attr)
			if !found && e.required {
				return nil, fmt.Errorf("Attribute '%s' for '%s' not found", e.attr, e.selector)
			}
			// TODO: Fail if == "" && required
			if attr == "" {
				return nil, nil
			}
			return strings.TrimSpace(attr), nil
		}
	}
	if sel.Length() > 1 && !e.multiple {
		return nil, fmt.Errorf("Multiple '%s' found (%d)", e.selector, sel.Length())
	}

	var err error
	ret := sel.Map(func(idx int, s *goquery.Selection) string {
		attr, found := s.Attr(e.attr)
		if !found {
			if e.required {
				err = fmt.Errorf("Attribute '%s' for '%s'[%d] not found", e.attr, e.selector, idx)
				return ""
			}
		}
		return strings.TrimSpace(attr)
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}
