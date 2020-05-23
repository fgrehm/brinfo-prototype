package core

import (
	"time"
)

type ContentSource struct {
	ID             string
	Host           string
	ArticleScraper ArticleScraper
}

type ContentSourceRepo interface {
	// TODO: Add context argument
	Register(cs *ContentSource) error
	FindByID(id string) (*ContentSource, error)
	FindByHost(id string) (*ContentSource, error)
}

type ArticleScraper interface {
	// TODO: Add context argument
	Run(articleHtml []byte, url string) (*ScrapedArticleData, error)
}

type ScrapedArticleData struct {
	Extra        interface{}            `json:"brinfo"`
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
