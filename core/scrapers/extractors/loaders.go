package extractors

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var extractorSpecRegexp = regexp.MustCompile(`^\s*([^\\|]+\S?)\s*\|\s*([\w]+)(\?)?(?:::(time))?\s*$`)

func FromString(extractorStr string) (Extractor, error) {
	match := extractorSpecRegexp.FindStringSubmatch(extractorStr)
	if len(match) < 3 {
		return nil, fmt.Errorf("Invalid extractor provided: %s", extractorStr)
	}

	selector := match[1]
	attribute := match[2]
	modifier := match[3]
	castTo := match[4]

	if castTo != "" && castTo != "time" {
		return nil, fmt.Errorf("cast to %s not supported", castTo)
	}

	required := true
	if modifier != "" {
		if modifier == "?" {
			required = false
		} else {
			return nil, fmt.Errorf("modifier %s not supported", modifier)
		}
	}

	if castTo == "time" {
		if attribute == "text" {
			if required {
				return TimeText(selector), nil
			} else {
				return OptTimeText(selector), nil
			}
		} else {
			if required {
				return TimeAttribute(selector, attribute), nil
			} else {
				return OptTimeAttribute(selector, attribute), nil
			}
		}
	} else if attribute == "text" {
		if required {
			return Text(selector, false), nil
		} else {
			return OptText(selector, false), nil
		}
	} else {
		if required {
			return Attribute(selector, attribute), nil
		} else {
			return OptAttribute(selector, attribute), nil
		}
	}
}

func FromJSON(data []byte) ([]Extractor, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		return nil, err
	}

	structuredExtractors := []Extractor{}
	htmlDocumentExtractors := map[string]Extractor{}
	for key, value := range jsonData {
		strValue, ok := value.(string)
		if ok {
			e, err := FromString(strValue)
			if err != nil {
				return nil, err
			}
			htmlDocumentExtractors[normalizeAttributeName(key)] = e
		} else {
			newExtractor, err := structuredFromMap(key, value.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			structuredExtractors = append(structuredExtractors, newExtractor)
		}
	}

	if len(htmlDocumentExtractors) > 0 {
		structuredExtractors = append(structuredExtractors, Structured("html", htmlDocumentExtractors))
	}

	return structuredExtractors, nil
}

func structuredFromMap(selector string, mapFromJSON map[string]interface{}) (Extractor, error) {
	extractors := map[string]Extractor{}
	for field, extractorStr := range mapFromJSON {
		newExtractor, err := FromString(extractorStr.(string))
		if err != nil {
			return nil, err
		}
		extractors[normalizeAttributeName(field)] = newExtractor
	}
	return Structured(selector, extractors), nil
}

func normalizeAttributeName(attr string) string {
	switch attr {
	case "published_at":
		return "publishedAt"
	case "updated_at", "modified_at", "modifiedAt":
		return "updatedAt"
	case "image_url":
		return "imageURL"
	case "full_text":
		return "fullText"
	default:
		return attr
	}
}
