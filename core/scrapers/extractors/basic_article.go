package extractors

import (
	"time"
)

type basicArticleExtractor struct{}

func BasicArticle() Extractor {
	return &basicArticleExtractor{}
}

func (e *basicArticleExtractor) Extract(args ExtractorArgs) (ExtractorResult, error) {
	result, err := HTMLInfo().Extract(args)
	if err != nil {
		return nil, err
	}

	data, ok := result.(map[string]interface{})
	if !ok {
		panic("Something unexpected returned from htmlinfo")
	}
	if data["publishedAt"] == (*time.Time)(nil) {
		if err = e.publishedAtFallbacks(data, args); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (e *basicArticleExtractor) publishedAtFallbacks(data map[string]interface{}, args ExtractorArgs) error {
	if data["publishedAt"] == nil && data["modifiedAt"] != nil {
		data["publishedAt"] = data["modifiedAt"]
		return nil
	}

	val, err := PublishedDates().Extract(args)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	}
	extractedData, ok := val.(map[string]*time.Time)
	if !ok {
		panic("Returned something weird")
	}
	data["publishedAt"] = extractedData["publishedAt"]
	data["modifiedAt"] = extractedData["modifiedAt"]
	// TODO: Annotate with the source of the date here

	return nil
}
