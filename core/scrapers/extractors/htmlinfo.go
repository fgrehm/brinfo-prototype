package extractors

import (
	"bytes"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dyatlov/go-htmlinfo/htmlinfo"
	"github.com/dyatlov/go-oembed/oembed"
)

var trimTitleSuffixRegexp = regexp.MustCompile(`(.+)\|.*$`)

type htmlinfoExtractor struct{}

func HTMLInfo() Extractor {
	return &htmlinfoExtractor{}
}

func (e *htmlinfoExtractor) Extract(args ExtractorArgs) (ExtractorResult, error) {
	html, err := args.Root.Html()
	if err != nil {
		return nil, err
	}

	htmlInfo := htmlinfo.NewHTMLInfo()
	if err := htmlInfo.Parse(bytes.NewBuffer([]byte(html)), &args.URL, &args.HTTPContentType); err != nil {
		return nil, err
	}

	return e.extractData(htmlInfo, args.URL)
}

func (e *htmlinfoExtractor) extractData(htmlInfo *htmlinfo.HTMLInfo, url string) (ExtractorResult, error) {
	oembed := htmlInfo.GenerateOembedFor(url)

	siteName := e.extractSiteName(htmlInfo, oembed)
	title := e.extractTitle(htmlInfo, oembed, siteName)
	excerpt := e.extractExcerpt(htmlInfo, oembed, title)
	fullText := e.extractFullText(htmlInfo, title, excerpt)
	imageURL := e.extractImageURL(htmlInfo, oembed)
	publishedAt, modifiedAt := e.extractDates(htmlInfo)

	// TODO: Consider returning canonical & opengraph url
	result := map[string]interface{}{
		"extra": map[string]interface{}{
			"htmlinfo": map[string]interface{}{
				"data":             htmlInfo,
				"generated_oembed": oembed,
			},
		},
		"title":       title,
		"excerpt":     excerpt,
		"fullText":    fullText,
		"imageURL":    imageURL,
		"publishedAt": publishedAt,
		"modifiedAt":  modifiedAt,
	}
	return result, nil
}

func (e htmlinfoExtractor) extractSiteName(info *htmlinfo.HTMLInfo, oembed *oembed.Info) string {
	if info.OGInfo != nil && info.OGInfo.SiteName != "" {
		return info.OGInfo.SiteName
	}
	if oembed != nil && oembed.ProviderName != "" {
		return oembed.ProviderName
	}
	return ""
}

func (e *htmlinfoExtractor) extractTitle(htmlInfo *htmlinfo.HTMLInfo, oembed *oembed.Info, siteName string) string {
	var title string
	if oembed != nil && oembed.Title != "" {
		title = oembed.Title
	}
	if title == "" && htmlInfo.Title != "" {
		title = htmlInfo.Title
	}
	if title != "" && siteName != "" {
		title = strings.TrimSuffix(title, " - "+siteName)
		title = strings.TrimSuffix(title, " | "+siteName)
	}
	if title != "" && trimTitleSuffixRegexp.MatchString(title) {
		title = trimTitleSuffixRegexp.ReplaceAllString(title, "${1}")
	}
	return strings.TrimSpace(title)
}

func (e *htmlinfoExtractor) extractExcerpt(htmlInfo *htmlinfo.HTMLInfo, oembed *oembed.Info, title string) string {
	var excerpt string
	if oembed != nil && oembed.Description != "" {
		excerpt = oembed.Description
	}
	if htmlInfo.Description != "" {
		excerpt = htmlInfo.Description
	}

	excerpt = strings.TrimSpace(excerpt)
	if !e.isGoodExcerpt(excerpt) {
		return ""
	}

	// TODO: Remove timestamps, like the ones in
	// http://www.amazonas.am.gov.br/2020/05/casa-do-migrante-jacamim-27-anos-acolhendo-pessoas-em-situacao-de-vulnerabilidade/
	return strings.TrimSpace(strings.TrimPrefix(excerpt, title))
}

func (e *htmlinfoExtractor) isGoodExcerpt(excerpt string) bool {
	if excerpt == "" || excerpt == "..." {
		return false
	}

	words := strings.Split(strings.Trim(excerpt, "..."), " ")
	return len(words) > 4
}

func (e *htmlinfoExtractor) extractFullText(htmlInfo *htmlinfo.HTMLInfo, title, excerpt string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlInfo.MainContent))
	if err != nil {
		panic(err)
	}

	fullText := doc.Find("body").Text()
	fullTextChunks := []string{}
	for _, str := range strings.Split(fullText, "\n") {
		str = strings.TrimSpace(str)
		if str != "" && str != title && str != excerpt {
			fullTextChunks = append(fullTextChunks, str)
		}
	}
	fullText = strings.Join(fullTextChunks, "\n")
	return fullText
}

func (e *htmlinfoExtractor) extractImageURL(htmlInfo *htmlinfo.HTMLInfo, oembed *oembed.Info) string {
	if htmlInfo.OGInfo != nil && len(htmlInfo.OGInfo.Images) > 0 {
		for _, img := range htmlInfo.OGInfo.Images {
			if img.Width != 0 && img.Height != 0 {
				if img.SecureURL != "" {
					return img.SecureURL
				}
				if img.URL != "" {
					return img.URL
				}
			}
		}
	}

	if oembed != nil {
		return oembed.ThumbnailURL
	}
	return ""
}

func (e *htmlinfoExtractor) extractDates(htmlInfo *htmlinfo.HTMLInfo) (*time.Time, *time.Time) {
	if htmlInfo.OGInfo == nil || htmlInfo.OGInfo.Article == nil {
		return nil, nil
	}

	publishedAt := htmlInfo.OGInfo.Article.PublishedTime
	modifiedAt := htmlInfo.OGInfo.Article.ModifiedTime
	return publishedAt, modifiedAt
}
