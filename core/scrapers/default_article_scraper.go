package scrapers

import (
	"bytes"
	"regexp"
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

	selection := doc.Selection

	s.fallbackPublishedAtFromMeta(data, selection)
	if data.PublishedAt == nil {
		err = s.fallbackPublishedAtFromRDF(data, selection)
	}

	return err
}

func (s *defaultArticleScraper) fallbackPublishedAtFromMeta(data *core.ScrapedArticleData, root *goquery.Selection) error {
	extractor := xt.Structured("head", map[string]xt.Extractor{
		"published_at": xt.OptAttribute(`meta[property="article:published_time"]`, "content"),
		"modified_at":  xt.OptAttribute(`meta[property="article:modified_time"]`, "content"),
	})

	extracted, err := extractor.Extract(root)
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

	if extractedMap["published_at"] != nil {
		publishedAt, err := s.parseDate(extractedMap["published_at"], brLoc)
		if err != nil {
			return err
		}
		data.PublishedAt = &publishedAt
	}

	if extractedMap["modified_at"] != nil {
		modifiedAt, err := s.parseDate(extractedMap["modified_at"], brLoc)
		if err != nil {
			return err
		}
		data.ModifiedAt = &modifiedAt
	}

	return nil
}

var (
	brDateTimeRegex = regexp.MustCompile(`^[0-9]{1,2}/[0-9]{1,2}/[0-9]{2,4}\s+[0-9]{1,2}h[0-9]{1,2}$`)
	defaultTime = time.Time{}
)

func (s *defaultArticleScraper) fallbackPublishedAtFromRDF(data *core.ScrapedArticleData, root *goquery.Selection) error {
	extractor := xt.Structured(`article[vocab*="schema.org"][typeof=Article][prefix*=rnews]`, map[string]xt.Extractor{
		"published_at": xt.OptText(`[property="rnews:datePublished"]`, false),
		"modified_at":  xt.OptText(`[property="rnews:dateModified"]`, false),
	})

	extracted, err := extractor.Extract(root)
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
	time.Local = brLoc

	if extractedMap["published_at"] != nil {
		publishedAt, err := s.parseDate(extractedMap["published_at"], brLoc)
		if err != nil {
			return err
		}
		if publishedAt != defaultTime {
			data.PublishedAt = &publishedAt
		}
	}

	if extractedMap["modified_at"] != nil {
		modifiedAt, err := s.parseDate(extractedMap["modified_at"], brLoc)
		if err != nil {
			return err
		}
		if modifiedAt != defaultTime {
			data.ModifiedAt = &modifiedAt
		}
	}

	return nil
}

func (*defaultArticleScraper) parseDate(datetime xt.ExtractorResult, loc *time.Location) (time.Time, error) {
	dateStr, ok := datetime.(string)
	if !ok {
		panic("Tried to parse something that is not a string")
	}

	var (
		dt  time.Time
		err error
	)

	if brDateTimeRegex.MatchString(dateStr) {
		dt, err = time.ParseInLocation("_2/01/2006 15h04", dateStr, loc)
		if err != nil {
			return time.Time{}, err
		}
	}

	if dt == defaultTime {
		dt, err = dateparse.ParseIn(dateStr, loc)
		if err != nil {
			return time.Time{}, err
		}
	}
	return dt, nil
}
