package core

import (
	"context"
	"encoding/json"
	"time"
)

type ArticleScraper interface {
	Run(ctx context.Context, html []byte, url, httpContentType string) (*ArticleData, error)
}

type ArticleListScraper interface {
	Run(ctx context.Context, html []byte, url, httpContentType string) ([]*ArticleLink, error)
}

type ArticleData struct {
	Extra        map[string]interface{} `json:"brinfo"`
	URL          string                 `json:"url"`
	URLHash      string                 `json:"url_hash"`
	Title        string                 `json:"title"`
	FullText     string                 `json:"full_text"`
	FullTextHash string                 `json:"full_text_hash"`
	Excerpt      string                 `json:"excerpt"`
	FoundAt      time.Time              `json:"found_at"`
	PublishedAt  *time.Time             `json:"published_at"`
	ModifiedAt   *time.Time             `json:"updated_at"` // TODO: Serialize to modified at after changing covid19br.pub
	ImageURL     string                 `json:"image_url"`
}

type ArticleLink struct {
	URL         string     `json:"url"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
}

func ArticleDataFromJSON(data []byte) (*ArticleData, error) {
	articleData := &ArticleData{}
	err := json.Unmarshal(data, &articleData)
	if err != nil {
		return nil, err
	}
	return articleData, nil
}

func (d *ArticleData) CollectValues(other *ArticleData) {
	if other.Extra != nil && len(other.Extra) > 0 {
		if d.Extra == nil || len(d.Extra) == 0 {
			d.Extra = other.Extra
		} else {
			for k, v := range other.Extra {
				d.Extra[k] = v
			}
		}
	}
	if other.URL != "" {
		d.URL = other.URL
		d.URLHash = other.URLHash
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
	// TODO: Make sure the modified at is >= pubat
	if other.PublishedAt != nil && !other.PublishedAt.IsZero() {
		d.PublishedAt = other.PublishedAt
	}
	if other.ModifiedAt != nil && !other.ModifiedAt.IsZero() {
		d.ModifiedAt = other.ModifiedAt
	}
	if other.ImageURL != "" {
		d.ImageURL = other.ImageURL
	}
}

func (d *ArticleData) ValidForIngestion() bool {
	valid := d.URL != ""
	valid = valid && d.URLHash != ""
	valid = valid && d.Title != ""
	valid = valid && d.FullText != ""
	valid = valid && d.FullTextHash != ""
	valid = valid && d.PublishedAt != nil && !d.PublishedAt.IsZero()
	valid = valid && (d.ModifiedAt == nil || !d.ModifiedAt.IsZero())
	valid = valid && !d.FoundAt.IsZero()
	return valid
}
