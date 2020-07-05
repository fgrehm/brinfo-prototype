package extractors

import (
	"time"
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

func (e *publishedDatesExtractor) Extract(args ExtractorArgs) (ExtractorResult, error) {
	result, err := e.extractWithFallbacks(args)
	if err != nil {
		return nil, err
	}

	if result == nil || (result.publishedAt == nil && result.modifiedAt == nil) {
		return nil, nil
	}

	return map[string]*time.Time{
		"publishedAt": result.publishedAt,
		"modifiedAt":  result.modifiedAt,
	}, nil
}

func (e *publishedDatesExtractor) extractWithFallbacks(args ExtractorArgs) (*extractedDates, error) {
	result, err := e.extractFromMeta(args)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	result, err = e.extractFromRNews(args)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	result, err = e.extractFromArticleTime(args)
	if err != nil {
		return nil, err
	}
	if result != nil {
		return result, nil
	}

	return nil, nil
}

func (e *publishedDatesExtractor) extractFromMeta(args ExtractorArgs) (*extractedDates, error) {
	extractor := Structured("head", map[string]Extractor{
		"published_at": OptTimeAttribute(`meta[property="article:published_time"]`, "content"),
		"modified_at":  OptTimeAttribute(`meta[property="article:modified_time"]`, "content"),
	})
	return e.handleExtractedResult(extractor.Extract(args))
}

func (e *publishedDatesExtractor) extractFromRNews(args ExtractorArgs) (*extractedDates, error) {
	extractor := Structured(`body [vocab*="schema.org"][typeof=Article][prefix*=rnews]`, map[string]Extractor{
		"published_at": OptTimeText(`[property="rnews:datePublished"]`),
		"modified_at":  OptTimeText(`[property="rnews:dateModified"]`),
	})

	return e.handleExtractedResult(extractor.Extract(args))
}

func (e *publishedDatesExtractor) extractFromArticleTime(args ExtractorArgs) (*extractedDates, error) {
	extractor := Structured(`article`, map[string]Extractor{
		"published_at": OptTimeAttribute(`time`, "pubdate"),
	})

	dates, err := e.handleExtractedResult(extractor.Extract(args))
	if err != nil {
		return nil, err
	}
	if dates != nil {
		return dates, nil
	}

	extractor = Structured(`article`, map[string]Extractor{
		"published_at": OptTimeAttribute(`time[pubdate]`, "datetime"),
	})
	return e.handleExtractedResult(extractor.Extract(args))
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
