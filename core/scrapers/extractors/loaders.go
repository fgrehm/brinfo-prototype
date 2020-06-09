package extractors

import (
	"fmt"
	"regexp"
)

var extractorSpecRegexp = regexp.MustCompile(`^\s*([^\\|]+\S)\s*\|\s*([\w]+)(\?)?(?:::(time))?\s*$`)

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
