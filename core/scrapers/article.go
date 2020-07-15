package scrapers

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"time"

	"github.com/fgrehm/brinfo/core"
	xt "github.com/fgrehm/brinfo/core/scrapers/extractors"

	"github.com/PuerkitoBio/goquery"
	"github.com/mitchellh/mapstructure"
)

type articleScraper struct {
	*ArticleScraperConfig
}

type ArticleScraperConfig struct {
	Clock      Clock
	Extractors []xt.Extractor
	MergeWith  *core.ArticleData
}

type Clock interface {
	Now() time.Time
}

func NewArticleScraper(cfg *ArticleScraperConfig) core.ArticleScraper {
	return &articleScraper{cfg}
}

func (s *articleScraper) Run(ctx context.Context, html []byte, url, httpContentType string) (*core.ArticleData, error) {
	data := &core.ArticleData{
		URL:     url,
		FoundAt: s.Clock.Now(),
		Extra: map[string]interface{}{
			"html": mustGzip(html),
		},
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(html))
	if err != nil {
		return nil, err
	}

	args := xt.ExtractorArgs{
		Context:         ctx,
		URL:             url,
		HTTPContentType: httpContentType,
		Root:            doc.Selection,
	}
	for _, extractor := range s.Extractors {
		result, err := extractor.Extract(args)
		if err != nil {
			return nil, err
		}

		extractorData := &core.ArticleData{}
		// TODO: Should error if something comes back that can't be mapped into the struct
		if err = mapstructure.Decode(result, extractorData); err != nil {
			return nil, err
		}
		data.CollectValues(extractorData)
	}

	// TODO: Test this
	if s.MergeWith != nil {
		data.CollectValues(s.MergeWith)
	}

	if data.URL != "" {
		data.URLHash = s.generateHash(data.URL)
	}

	if data.FullText != "" {
		data.FullTextHash = s.generateHash(data.FullTextHash)
	}

	if data.ModifiedAt != nil && data.PublishedAt == nil {
		data.PublishedAt = data.ModifiedAt
	}

	return data, nil
}

func (s *articleScraper) generateHash(text string) string {
	algorithm := sha1.New()
	if _, err := algorithm.Write([]byte(text)); err != nil {
		panic(err)
	}
	return hex.EncodeToString(algorithm.Sum(nil))
}
