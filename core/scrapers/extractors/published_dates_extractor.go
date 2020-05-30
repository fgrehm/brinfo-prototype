package extractors

import (
	"time"

	"github.com/PuerkitoBio/goquery"
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
		"modified_at":  result.modifiedAt,
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
		"published_at": OptTimeAttribute(`meta[property="article:published_time"]`, "content"),
		"modified_at":  OptTimeAttribute(`meta[property="article:modified_time"]`, "content"),
	})
	return e.handleExtractedResult(extractor.Extract(root))
}

func (e *publishedDatesExtractor) extractFromRNews(root *goquery.Selection) (*extractedDates, error) {
	extractor := Structured(`body [vocab*="schema.org"][typeof=Article][prefix*=rnews]`, map[string]Extractor{
		"published_at": OptTimeText(`[property="rnews:datePublished"]`),
		"modified_at":  OptTimeText(`[property="rnews:dateModified"]`),
	})

	return e.handleExtractedResult(extractor.Extract(root))
}

func (e *publishedDatesExtractor) extractFromArticleTime(root *goquery.Selection) (*extractedDates, error) {
	extractor := Structured(`article`, map[string]Extractor{
		"published_at": OptTimeAttribute(`time`, "pubdate"),
	})

	dates, err := e.handleExtractedResult(extractor.Extract(root))
	if err != nil {
		return nil, err
	}
	if dates != nil {
		return dates, nil
	}

	extractor = Structured(`article`, map[string]Extractor{
		"published_at": OptTimeAttribute(`time[pubdate]`, "datetime"),
	})
	return e.handleExtractedResult(extractor.Extract(root))
}

func (e *publishedDatesExtractor) handleExtractedResult(extracted ExtractorResult, err error) (*extractedDates, error) {
	if extracted == nil {
		return nil, nil
	}

	extractedMap, ok := extracted.(map[string]ExtractorResult)
	if !ok {
		panic("Extractor returned something weird")
	}
	if len(extractedMap) == 0 {
		return nil, nil
	}
	if extractedMap["published_at"] == nil && extractedMap["modified_at"] == nil {
		return nil, nil
	}

	data := &extractedDates{}
	publishedAt := extractedMap["published_at"]
	if publishedAt != nil {
		pubAt := publishedAt.(time.Time)
		data.publishedAt = &pubAt
	}

	modifiedAt := extractedMap["modified_at"]
	if modifiedAt != nil {
		modAt := modifiedAt.(time.Time)
		data.modifiedAt = &modAt
	}

	return data, nil
}
