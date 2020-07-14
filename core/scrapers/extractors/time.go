package extractors

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

var (
	brDateTimeRegex            = regexp.MustCompile(`^[0-9]{1,2}/[0-9]{1,2}/[0-9]{2,4}(\s|-)+[0-9]{1,2}[h:][0-9]{1,2}([m:][0-9]{1,2}s?)?$`)
	brDateTimeWithSecondsRegex = regexp.MustCompile(`^[0-9]{1,2}/[0-9]{1,2}/[0-9]{2,4}\s+[0-9]{1,2}:[0-9]{1,2}:[0-9]{1,2}$`)

	brLongDateTimeRegex     = regexp.MustCompile(`(?i)(?P<day>\d{1,2}) de (?P<longMonth>(\w|ç)+) de (?P<longYear>\d{4}) [àa]s (?P<hour>\d{1,2}):(?P<minute>\d{1,2})`)
	brLongDateTimeAmPmRegex = regexp.MustCompile(`(?i)(?P<day>\d{1,2})/(?P<longMonth>(\w|ç)+)/(?P<longYear>\d{4}) (?P<hour>\d{1,2}):(?P<minute>\d{1,2})\s*(?P<ampm>am|pm)`)

	brLoc *time.Location
)

type timeAttrExtractor struct {
	*attrExtractor
}

type timeTextExtractor struct {
	*textExtractor
	format string
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

func (e *timeAttrExtractor) Extract(args ExtractorArgs) (ExtractorResult, error) {
	res, err := e.attrExtractor.Extract(args)
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

func (e *timeTextExtractor) Extract(args ExtractorArgs) (ExtractorResult, error) {
	res, err := e.textExtractor.Extract(args)
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

func parseExtractedTime(timeStr string) (*time.Time, error) {
	var (
		dt  time.Time
		err error
	)

	if brDateTimeRegex.MatchString(timeStr) {
		timeStr = strings.Replace(timeStr, "h", ":", 1)
		timeStr = strings.Replace(timeStr, "H", ":", 1)
		timeStr = strings.Replace(timeStr, "m", ":", 1)
		timeStr = strings.Replace(timeStr, "M", ":", 1)
		timeStr = strings.Replace(timeStr, " - ", " ", 1)
		// HACK
		if brDateTimeWithSecondsRegex.MatchString(timeStr) {
			dt, err = time.ParseInLocation("_2/01/2006 15:04:05", timeStr, brLoc)
		} else {
			dt, err = time.ParseInLocation("_2/01/2006 15:04", timeStr, brLoc)
		}
		if err != nil {
			return nil, err
		}
	} else if brLongDateTimeRegex.MatchString(timeStr) {
		timeStr = adjustLongDateTime(timeStr, brLongDateTimeRegex)
	} else if brLongDateTimeAmPmRegex.MatchString(timeStr) {
		timeStr = adjustLongDateTime(timeStr, brLongDateTimeAmPmRegex)
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

var longMonthTranslationMap = map[*regexp.Regexp]string{
	regexp.MustCompile(`(?i)janeiro`):   "01",
	regexp.MustCompile(`(?i)fevereiro`): "02",
	regexp.MustCompile(`(?i)março`):     "03",
	regexp.MustCompile(`(?i)marco`):     "03",
	regexp.MustCompile(`(?i)abril`):     "04",
	regexp.MustCompile(`(?i)maio`):      "05",
	regexp.MustCompile(`(?i)junho`):     "06",
	regexp.MustCompile(`(?i)julho`):     "07",
	regexp.MustCompile(`(?i)agosto`):    "08",
	regexp.MustCompile(`(?i)setembro`):  "09",
	regexp.MustCompile(`(?i)outubro`):   "10",
	regexp.MustCompile(`(?i)novembro`):  "11",
	regexp.MustCompile(`(?i)dezembro`):  "11",
}

func adjustLongDateTime(timeStr string, formatRegex *regexp.Regexp) string {
	template := []byte("$month$longMonth/$day/$year$longYear $hour:$minute$ampm")
	result := []byte{}
	for _, submatches := range formatRegex.FindAllSubmatchIndex([]byte(timeStr), -1) {
		result = formatRegex.Expand(result, template, []byte(timeStr), submatches)
	}
	return translateLongMonth(string(result))
}

func translateLongMonth(str string) string {
	for reg, replacement := range longMonthTranslationMap {
		if reg.MatchString(str) {
			return reg.ReplaceAllString(str, replacement)
		}
	}
	return str
}
