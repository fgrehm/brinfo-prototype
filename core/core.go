package core

import (
	"context"
	"encoding/json"
	"net/url"
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

func (d *ArticleData) ValidForIngestion() (bool, []string) {
	errors := []string{}
	// now := time.Now()

	if d.URL == "" {
		errors = append(errors, "missing url")
	}
	if d.URLHash == "" {
		errors = append(errors, "missing url_hash")
	}
	if d.Title == "" {
		errors = append(errors, "missing title")
	}
	if d.FullText == "" {
		errors = append(errors, "missing full_text")
	}
	if d.FullTextHash == "" {
		errors = append(errors, "missing full_text_hash")
	}
	if d.PublishedAt == nil || d.PublishedAt.IsZero() {
		errors = append(errors, "missing published_at")
	}
	if d.ModifiedAt != nil && d.ModifiedAt.IsZero() {
		errors = append(errors, "updated_at in the future")
	}

	now := time.Now()
	if d.PublishedAt != nil && d.PublishedAt.Sub(now) >= (time.Hour*12) {
		errors = append(errors, "published_at in the future")
	}
	if d.ModifiedAt != nil && d.ModifiedAt.Sub(now) >= (time.Hour*12) {
		errors = append(errors, "updated_at in the future")
	}

	if d.FoundAt.IsZero() {
		errors = append(errors, "missing found_at")
	}
	if d.ImageURL != "" {
		u, err := url.Parse(d.ImageURL)
		if err != nil {
			panic(err)
		}
		if u.Host == "" {
			errors = append(errors, "image_url is not absolute")
		}
	}

	return len(errors) == 0, errors
}
