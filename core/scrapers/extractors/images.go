package extractors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fgrehm/brinfo/core"
)

type imgExtractor struct {
	selector string
	attr     string
}

func Images(selector, attr string) Extractor {
	return &imgExtractor{
		selector: selector,
		attr:     attr,
	}
}

func (e *imgExtractor) Extract(root *goquery.Selection) (ExtractorResult, error) {
	sel := root.Find(e.selector)
	if sel.Length() == 0 {
		return nil, nil
	}

	var err error
	ret := []*core.ScrapedArticleImage{}
	sel.Each(func(idx int, s *goquery.Selection) {
		url := e.getAttr(s, e.attr)
		if url == "" {
			err = fmt.Errorf("Attribute '%s' for '%s'[%d] not found", e.attr, e.selector, idx)
			return
		}

		imgData := &core.ScrapedArticleImage{Url: url}

		if widthStr := e.getAttr(s, "width"); widthStr != "" {
			width, err := strconv.ParseUint(widthStr, 10, 32)
			if err == nil {
				imgData.Width = width
			}
		}
		if heightStr := e.getAttr(s, "height"); heightStr != "" {
			height, err := strconv.ParseUint(heightStr, 10, 32)
			if err == nil {
				imgData.Height = height
			}
		}

		ret = append(ret, imgData)
	})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (e *imgExtractor) getAttr(sel *goquery.Selection, attr string) string {
	val, found := sel.Attr(attr)
	if !found {
		return ""
	}
	return strings.Trim(val, " ")
}
