package extractors

import (
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
)

type publishedDatesExtractor struct {
	brLoc time.Location
}

type extractedDates struct {
	publishedAt *time.Time
	modifiedAt  *time.Time
}

func PublishedDates() Extractor {
	brLoc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	return &publishedDatesExtractor{brLoc: *brLoc}
}

var brDateTimeRegex = regexp.MustCompile(`^[0-9]{1,2}/[0-9]{1,2}/[0-9]{2,4}\s+[0-9]{1,2}h[0-9]{1,2}$`)

func (e *publishedDatesExtractor) Extract(root *goquery.Selection) (ExtractorResult, error) {
	result, err := e.extractWithFallbacks(root)
	if err != nil {
		return nil, err
	}

	if result == nil || (result.publishedAt == nil && result.modifiedAt == nil) {
		return nil, nil
	}

	return map[string]*time.Time{
		"published_at": result.publishedAt,
		"modified_at": result.modifiedAt,
	}, nil
}

func (e *publishedDatesExtractor) extractWithFallbacks(root *goquery.Selection) (*extractedDates, error) {
	result, err := e.extractFromMeta(root)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	result, err = e.extractFromRNews(root)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	result, err = e.extractFromArticleTime(root)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	return nil, nil
}

func (e *publishedDatesExtractor) extractFromMeta(root *goquery.Selection) (*extractedDates, error) {
	extractor := Structured("head", map[string]Extractor{
		"published_at": OptAttribute(`meta[property="article:published_time"]`, "content"),
		"modified_at":  OptAttribute(`meta[property="article:modified_time"]`, "content"),
	})
	return e.handleExtractedResult(extractor.Extract(root))
}

func (e *publishedDatesExtractor) extractFromRNews(root *goquery.Selection) (*extractedDates, error) {
	extractor := Structured(`body [vocab*="schema.org"][typeof=Article][prefix*=rnews]`, map[string]Extractor{
		"published_at": OptText(`[property="rnews:datePublished"]`, false),
		"modified_at":  OptText(`[property="rnews:dateModified"]`, false),
	})

	return e.handleExtractedResult(extractor.Extract(root))
}

func (e *publishedDatesExtractor) extractFromArticleTime(root *goquery.Selection) (*extractedDates, error) {
	extractor := Structured(`article`, map[string]Extractor{
		"published_at": OptAttribute(`time[pubdate]`, "datetime"),
	})

	return e.handleExtractedResult(extractor.Extract(root))
}

func (e *publishedDatesExtractor) handleExtractedResult(extracted ExtractorResult, err error) (*extractedDates, error) {
	extractedMap, ok := extracted.(map[string]ExtractorResult)
	if !ok {
		panic("Extractor returned something weird")
	}
	if extractedMap["published_at"] == nil && extractedMap["modified_at"] == nil {
		return nil, nil
	}

	data := &extractedDates{}
	publishedAt, err := e.parseDate(extractedMap, "published_at")
	if err != nil {
		return nil, err
	}
	if publishedAt != nil {
		data.publishedAt = publishedAt
	}

	modifiedAt, err := e.parseDate(extractedMap, "modified_at")
	if err != nil {
		return nil, err
	}
	if modifiedAt != nil {
		data.modifiedAt = modifiedAt
	}

	return data, nil
}

func (e *publishedDatesExtractor) parseDate(extractedMap map[string]ExtractorResult, field string) (*time.Time, error) {
	datetime := extractedMap[field]
	if datetime == nil {
		return nil, nil
	}

	dateStr, ok := datetime.(string)
	if !ok {
		panic("Tried to parse something that is not a string")
	}

	var (
		dt  time.Time
		err error
	)

	if brDateTimeRegex.MatchString(dateStr) {
		dt, err = time.ParseInLocation("_2/01/2006 15h04", dateStr, &e.brLoc)
		if err != nil {
			return nil, err
		}
	}

	if dt.IsZero() {
		dt, err = dateparse.ParseIn(dateStr, &e.brLoc)
		if err != nil {
			return nil, err
		}
	}

	if dt.IsZero() {
		return nil, nil
	}
	return &dt, nil
}
