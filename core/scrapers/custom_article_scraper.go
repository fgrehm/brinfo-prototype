package scrapers

import (
	"bytes"
	"errors"
	"time"

	"github.com/fgrehm/brinfo/core"
	xt "github.com/fgrehm/brinfo/core/scrapers/extractors"

	"github.com/PuerkitoBio/goquery"
)

type customArticleScraper struct {
	cfg CustomArticleScraperConfig
}

type CustomArticleScraperConfig struct {
	PublishedAt xt.Extractor
	Images      xt.Extractor
}

func CustomArticleScraper(cfg CustomArticleScraperConfig) core.ArticleScraper {
	return &customArticleScraper{cfg}
}

func (s *customArticleScraper) Run(articleHtml []byte, url, contentType string) (*core.ScrapedArticleData, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(articleHtml))
	if err != nil {
		return nil, err
	}

	data := &core.ScrapedArticleData{}
	if err = s.extractPublishedAt(data, doc.Selection); err != nil {
		return nil, err
	}
	// TODO: Pass in URL for fixing relative links
	if err = s.extractImages(data, doc.Selection); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *customArticleScraper) extractPublishedAt(data *core.ScrapedArticleData, root *goquery.Selection) error {
	if s.cfg.PublishedAt == nil {
		return nil
	}

	val, err := s.cfg.PublishedAt.Extract(root)
	if err != nil {
		return err
	}

	pubAt, ok := val.(time.Time)
	if !ok {
		return errors.New("invalid type for for PublishedAt")
	}

	data.PublishedAt = &pubAt
	return nil
}

func (s *customArticleScraper) extractImages(data *core.ScrapedArticleData, root *goquery.Selection) error {
	if s.cfg.Images == nil {
		return nil
	}

	panic("Not working yet")

	// val, err := s.cfg.Images.Extract(root)
	// if err != nil {
	// 	return err
	// }

	// pubAt, ok := val.([]ExtractorResult)
	// if !ok {
	// 	return errors.New("invalid type for for Images")
	// }

	// data.PublishedAt = &pubAt
	return nil
}
