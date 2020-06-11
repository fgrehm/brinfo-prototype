package operations

import (
	"bytes"
	"context"
	neturl "net/url"
	"regexp"
	"time"

	"github.com/fgrehm/brinfo/core"
	xt "github.com/fgrehm/brinfo/core/scrapers/extractors"

	"github.com/PuerkitoBio/goquery"
)

type ScrapeArticlesListingArgs struct {
	UseCache             bool
	URL                  string
	LinkContainer        string
	URLExtractor         string
	PublishedAtExtractor string
	ImageURLExtractor    string
}

type articlesListingScraper struct {
	url       string
	parsedURL *neturl.URL
	extractor xt.Extractor
	cache     bool
}

func ScrapeArticlesListing(ctx context.Context, args ScrapeArticlesListingArgs) ([]*core.ArticleLink, error) {
	if args.LinkContainer == "" {
		panic("no link container")
	}
	if args.URLExtractor == "" {
		panic("no url extractor")
	}

	extractors := map[string]xt.Extractor{}

	e, err := xt.FromString(args.URLExtractor)
	if err != nil {
		return nil, err
	}
	extractors["url"] = e

	if args.PublishedAtExtractor != "" {
		e, err = xt.FromString(args.PublishedAtExtractor)
		if err != nil {
			return nil, err
		}
		extractors["published_at"] = e
	}

	if args.ImageURLExtractor != "" {
		e, err = xt.FromString(args.ImageURLExtractor)
		if err != nil {
			return nil, err
		}
		extractors["image_url"] = e
	}

	parsedURL, err := neturl.Parse(args.URL)
	if err != nil {
		panic(err)
	}

	scraper := &articlesListingScraper{
		url:       args.URL,
		parsedURL: parsedURL,
		extractor: xt.StructuredList(args.LinkContainer, extractors),
		cache:     args.UseCache,
	}

	return scraper.scrape(ctx)
}

func (s *articlesListingScraper) scrape(ctx context.Context) ([]*core.ArticleLink, error) {
	body, _, err := makeRequest(s.cache, s.url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	data, err := s.extractor.Extract(doc.Selection)
	if err != nil {
		return nil, err
	}

	list, ok := data.([]map[string]xt.ExtractorResult)
	if !ok {
		panic("Something went wrong")
	}

	ret := []*core.ArticleLink{}
	for _, res := range list {
		link := &core.ArticleLink{}

		if res["published_at"] != nil {
			pubAt := res["published_at"].(time.Time)
			link.PublishedAt = &pubAt
		}
		if res["image_url"] != nil {
			imageURL := s.fixRelativeURL(res["image_url"].(string))
			link.ImageURL = &imageURL
		}
		link.URL = s.fixRelativeURL(res["url"].(string))

		ret = append(ret, link)
	}

	return ret, nil
}

func (s *articlesListingScraper) fixRelativeURL(url string) string {
	url = regexp.MustCompile(`\s*`).ReplaceAllString(url, "")
	u, err := neturl.Parse(url)
	if err != nil {
		panic(err)
	}

	if u.Scheme == "" {
		u.Scheme = s.parsedURL.Scheme
	}
	if u.Host == "" {
		u.Host = s.parsedURL.Host
	}

	return u.String()
}
