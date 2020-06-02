package scrapers

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/fgrehm/brinfo/core"
	xt "github.com/fgrehm/brinfo/core/scrapers/extractors"

	"github.com/PuerkitoBio/goquery"
)

type customArticleScraper struct {
	cfg CustomArticleScraperConfig
}

type CustomArticleScraperConfig struct {
	Title       xt.Extractor
	PublishedAt xt.Extractor
	Images      xt.Extractor
}

func CustomArticleScraper(cfg CustomArticleScraperConfig) core.ArticleScraper {
	return &customArticleScraper{cfg}
}

func (s *customArticleScraper) Run(ctx context.Context, articleHtml []byte, url, contentType string) (*core.ScrapedArticleData, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(articleHtml))
	if err != nil {
		return nil, err
	}

	data := &core.ScrapedArticleData{}
	if err = s.extractTitle(data, doc.Selection); err != nil {
		return nil, err
	}
	if err = s.extractPublishedAt(data, doc.Selection); err != nil {
		return nil, err
	}
	if err = s.extractImages(data, doc.Selection, url); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *customArticleScraper) extractTitle(data *core.ScrapedArticleData, root *goquery.Selection) error {
	if s.cfg.Title == nil {
		return nil
	}

	val, err := s.cfg.Title.Extract(root)
	if err != nil {
		return err
	}

	title, ok := val.(string)
	if !ok {
		return errors.New("invalid type for for PublishedAt")
	}

	data.Title = title
	return nil
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

func (s *customArticleScraper) extractImages(data *core.ScrapedArticleData, root *goquery.Selection, url string) error {
	if s.cfg.Images == nil {
		return nil
	}

	val, err := s.cfg.Images.Extract(root)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	}

	imgs, ok := val.([]*core.ScrapedArticleImage)
	if !ok {
		return errors.New("invalid type for for Images")
	}
	s.fixRelativeImgUrls(imgs, url)

	data.Images = imgs
	if data.ImageUrl == "" {
		data.ImageUrl = imgs[0].Url
	}

	return nil
}

func (s *customArticleScraper) fixRelativeImgUrls(imgs []*core.ScrapedArticleImage, articleURL string) error {
	parsedArticleURL, err := url.Parse(articleURL)
	if err != nil {
		return err
	}

	for _, img := range imgs {
		u, err := url.Parse(img.Url)
		if err != nil {
			return err
		}

		if u.Scheme == "" {
			u.Scheme = parsedArticleURL.Scheme
		}
		if u.Host == "" {
			u.Host = parsedArticleURL.Host
		}

		img.Url = u.String()
	}
	return nil
}
