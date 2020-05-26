package scrapers

import (
	"bytes"
	"time"

	"github.com/fgrehm/brinfo/core"
	xt "github.com/fgrehm/brinfo/core/scrapers/extractors"

	"github.com/PuerkitoBio/goquery"
)

var DefaultArticleScraper core.ArticleScraper

func init() {
	DefaultArticleScraper = &defaultArticleScraper{}
}

type defaultArticleScraper struct{}

func (f *defaultArticleScraper) Run(articleHtml []byte, url, contentType string) (*core.ScrapedArticleData, error) {
	htmlinfo := &htmlInfoScraper{}
	data, err := htmlinfo.Run(articleHtml, url, contentType)
	if err != nil {
		return nil, err
	}

	if data.PublishedAt == nil {
		err = f.publishedAtFallbacks(data, articleHtml)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (s *defaultArticleScraper) publishedAtFallbacks(data *core.ScrapedArticleData, articleHtml []byte) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(articleHtml))
	if err != nil {
		return err
	}

	val, err := xt.PublishedDates().Extract(doc.Selection)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	}
	extractedData, ok := val.(map[string]*time.Time)
	if !ok {
		panic("Returned something weird")
	}
	if publishedAt := extractedData["published_at"]; publishedAt != nil {
		data.PublishedAt = publishedAt
	}
	if modifiedAt := extractedData["modified_at"]; modifiedAt != nil {
		data.PublishedAt = modifiedAt
		if data.PublishedAt == nil {
			data.PublishedAt = modifiedAt
		}
	}

	return err
}
