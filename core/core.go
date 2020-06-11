package core

import (
	"context"
	"encoding/json"
	"time"
)

type ContentSource struct {
	ID               string
	Host             string
	ArticleScraper   ArticleScraper
	ForceContentType string
}

type ContentSourceRepo interface {
	Register(cs *ContentSource) error
	FindByID(ctx context.Context, id string) (*ContentSource, error)
	GetByHost(ctx context.Context, host string) (*ContentSource, error)
	FindByHost(ctx context.Context, host string) (*ContentSource, error)
}

type ArticleScraper interface {
	Run(ctx context.Context, articleHtml []byte, url, contentType string) (*ScrapedArticleData, error)
}

type combinedArticleScraper struct {
	scrapers []ArticleScraper
}

func CombinedArticleScraper(scrapers ...ArticleScraper) ArticleScraper {
	if len(scrapers) < 2 {
		panic("Need at least 2 scrapers to combine")
	}
	return &combinedArticleScraper{scrapers}
}

func (s *combinedArticleScraper) Run(ctx context.Context, articleHtml []byte, url, contentType string) (*ScrapedArticleData, error) {
	data := &ScrapedArticleData{}
	for _, scraper := range s.scrapers {
		newData, err := scraper.Run(ctx, articleHtml, url, contentType)
		if err != nil {
			return nil, err
		}
		data.CollectValues(newData)
	}
	return data, nil
}

type ScrapedArticleData struct {
	Extra        map[string]interface{} `json:"brinfo"`
	SourceID     string                 `json:"source_guid"`
	ContentType  string                 `json:"content_type"`
	Url          string                 `json:"url"`
	UrlHash      string                 `json:"url_hash"`
	Title        string                 `json:"title"`
	FullText     string                 `json:"full_text"`
	FullTextHash string                 `json:"full_text_hash"`
	Excerpt      string                 `json:"excerpt"`
	FoundAt      time.Time              `json:"found_at"`
	PublishedAt  *time.Time             `json:"published_at"`
	ModifiedAt   *time.Time             `json:"updated_at"` // TODO: Serialize to modified at after changing covid19br.pub
	Images       []*ScrapedArticleImage `json:"images"`
	ImageUrl     string                 `json:"image_url"`
}

func ScrapedArticleDataFromJSON(data []byte) (*ScrapedArticleData, error) {
	articleData := &ScrapedArticleData{}
	err := json.Unmarshal(data, &articleData)
	if err != nil {
		return nil, err
	}
	return articleData, nil
}

func (d *ScrapedArticleData) CollectValues(other *ScrapedArticleData) {
	if other.Extra != nil && len(other.Extra) > 0 {
		if d.Extra == nil || len(d.Extra) == 0 {
			d.Extra = other.Extra
		} else {
			for k, v := range other.Extra {
				d.Extra[k] = v
			}
		}
	}
	if other.SourceID != "" {
		d.SourceID = other.SourceID
	}
	if other.ContentType != "" {
		d.ContentType = other.ContentType
	}
	if other.Url != "" {
		d.Url = other.Url
		d.UrlHash = other.UrlHash
	}
	if other.Title != "" {
		d.Title = other.Title
	}
	if other.FullText != "" {
		d.FullText = other.FullText
		d.FullTextHash = other.FullTextHash
	}
	if other.Excerpt != "" {
		d.Excerpt = other.Excerpt
	}
	if !other.FoundAt.IsZero() {
		d.FoundAt = other.FoundAt
	}
	if other.PublishedAt != nil && !other.PublishedAt.IsZero() {
		d.PublishedAt = other.PublishedAt
	}
	if other.ModifiedAt != nil && !other.ModifiedAt.IsZero() {
		d.ModifiedAt = other.ModifiedAt
	}
	if len(other.Images) > 0 {
		d.Images = append(d.Images, other.Images...)
	}
	if other.ImageUrl != "" {
		d.ImageUrl = other.ImageUrl
	}
}

func (d *ScrapedArticleData) ValidForIngestion() bool {
	valid := d.SourceID != ""
	valid = valid && d.ContentType != ""
	valid = valid && d.Url != ""
	valid = valid && d.UrlHash != ""
	valid = valid && d.Title != ""
	valid = valid && d.FullText != ""
	valid = valid && d.FullTextHash != ""
	valid = valid && d.PublishedAt != nil && !d.PublishedAt.IsZero()
	valid = valid && (d.ModifiedAt == nil || !d.ModifiedAt.IsZero())
	valid = valid && !d.FoundAt.IsZero()
	return valid
}

type ScrapedArticleImage struct {
	Url       string `json:"url"`
	SecureUrl string `json:"secure_url"`
	Type      string `json:"type"`
	Width     uint64 `json:"width"`
	Height    uint64 `json:"height"`
}

type ArticleLink struct {
	URL         string     `json:"url"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
}
