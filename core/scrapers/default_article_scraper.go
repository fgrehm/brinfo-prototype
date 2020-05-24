package scrapers

import (
	"bytes"
	"time"

	"github.com/fgrehm/brinfo/core"
	xt "github.com/fgrehm/brinfo/core/scrapers/extractors"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
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
		err = f.fallbackPublishedAtFromMeta(data, articleHtml)
		// if err != nil {
		// 	return nil, err
		// }
	}

	return data, nil
}

func (s *defaultArticleScraper) fallbackPublishedAtFromMeta(data *core.ScrapedArticleData, articleHtml []byte) error {
	// Lookup `article:publishedat`
	extractor := xt.Structured("head", map[string]xt.Extractor{
		"published_at": xt.Attribute(`meta[property="article:published_time"]`, "content"),
		"modified_at":  xt.Attribute(`meta[property="article:modified_time"]`, "content"),
	})

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(articleHtml))
	if err != nil {
		return err
	}

	extracted, err := extractor.Extract(doc.Selection)
	if err != nil {
		return err
	}

	extractedMap, ok := extracted.(map[string]xt.ExtractorResult)
	if !ok {
		panic("Extractor returned something weird")
	}

	brLoc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	publishedAt, err := s.parseDate(extractedMap["published_at"], brLoc)
	if err != nil {
		return err
	}
	data.PublishedAt = &publishedAt

	modifiedAt, err := s.parseDate(extractedMap["modified_at"], brLoc)
	if err != nil {
		return err
	}
	data.ModifiedAt = &modifiedAt
	// TODO: merge data

	return nil
}

func (*defaultArticleScraper) parseDate(datetime xt.ExtractorResult, loc *time.Location) (time.Time, error) {
	dateStr, ok := datetime.(string)
	if !ok {
		panic("Tried to parse something that is not a string")
	}
	date, err := dateparse.ParseIn(dateStr, loc)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
