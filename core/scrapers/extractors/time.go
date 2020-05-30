package extractors

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
)

var (
	brDateTimeRegex = regexp.MustCompile(`^[0-9]{1,2}/[0-9]{1,2}/[0-9]{2,4}(\s|-)+[0-9]{1,2}[h:][0-9]{1,2}$`)
	brLoc           *time.Location
)

type timeAttrExtractor struct {
	*attrExtractor
}

type timeTextExtractor struct {
	*textExtractor
}

func init() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}
	brLoc = loc
}

func TimeAttribute(selector, attr string) Extractor {
	return &timeAttrExtractor{
		attrExtractor: &attrExtractor{
			selector: selector,
			attr:     attr,
			multiple: false,
			required: true,
		},
	}
}

func OptTimeAttribute(selector, attr string) Extractor {
	return &timeAttrExtractor{
		attrExtractor: &attrExtractor{
			selector: selector,
			attr:     attr,
			multiple: false,
			required: false,
		},
	}
}

func (e *timeAttrExtractor) Extract(root *goquery.Selection) (ExtractorResult, error) {
	res, err := e.attrExtractor.Extract(root)
	if err != nil {
		return nil, err
	}
	if res == nil {
		if e.attrExtractor.required {
			return nil, errors.New("unable to parse time")
		}
		return nil, nil
	}

	str, ok := res.(string)
	if !ok {
		return nil, errors.New("attrExtractor returned something not a string")
	}

	time, err := parseExtractedTime(str)
	if err != nil {
		return nil, err
	}
	if time != nil {
		return *time, nil
	}

	if e.attrExtractor.required {
		return nil, errors.New("unable to parse time")
	}

	return nil, nil
}

func (e *timeTextExtractor) Extract(root *goquery.Selection) (ExtractorResult, error) {
	res, err := e.textExtractor.Extract(root)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	str, ok := res.(string)
	if !ok {
		return nil, errors.New("attrExtractor returned something not a string")
	}

	time, err := parseExtractedTime(str)
	if err != nil {
		return nil, err
	}
	if time != nil {
		return *time, nil
	}

	if e.textExtractor.required {
		return nil, errors.New("unable to parse time")
	}

	return nil, nil
}

func TimeText(selector string) Extractor {
	return &timeTextExtractor{
		textExtractor: &textExtractor{
			selector: selector,
			multiple: false,
			required: true,
		},
	}
}

func OptTimeText(selector string) Extractor {
	return &timeTextExtractor{
		textExtractor: &textExtractor{
			selector: selector,
			multiple: false,
			required: false,
		},
	}
}

func parseExtractedTime(timeStr string) (*time.Time, error) {
	var (
		dt  time.Time
		err error
	)

	if brDateTimeRegex.MatchString(timeStr) {
		timeStr = strings.Replace(timeStr, "h", ":", 1)
		timeStr = strings.Replace(timeStr, "H", ":", 1)
		timeStr = strings.Replace(timeStr, " - ", " ", 1)
		dt, err = time.ParseInLocation("_2/01/2006 15:04", timeStr, brLoc)
		if err != nil {
			return nil, err
		}
	}

	if dt.IsZero() {
		dt, err = dateparse.ParseIn(timeStr, brLoc)
		if err != nil {
			return nil, err
		}
	}

	if dt.IsZero() {
		return nil, nil
	}
	return &dt, nil
}
