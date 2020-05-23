package scrapers

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"

	"github.com/fgrehm/brinfo/core"

	"github.com/PuerkitoBio/goquery"
	"github.com/dyatlov/go-htmlinfo/htmlinfo"
	"github.com/dyatlov/go-oembed/oembed"
	"github.com/dyatlov/go-opengraph/opengraph"
	// log "github.com/sirupsen/logrus"
)

var DefaultArticleScraper core.ArticleScraper

func init() {
	DefaultArticleScraper = &defaultArticleScraper{}
}

type defaultArticleScraper struct{}

func (f *defaultArticleScraper) Run(articleHtml []byte, url string) (*core.ScrapedArticleData, error) {
	info := htmlinfo.NewHTMLInfo()
	if err := info.Parse(bytes.NewBuffer(articleHtml), &url, nil); err != nil {
		return nil, err
	}
	oembed := info.GenerateOembedFor(url)

	data := &core.ScrapedArticleData{ContentType: "article"}

	data.Extra = map[string]interface{}{
		// TODO: Store full HTML
		// "html": buf.String(),
		"html_info": map[string]interface{}{
			"data":             info,
			"generated_oembed": oembed,
		},
	}

	siteName := getSiteName(info, oembed)
	title := cleanTitle(getTitleFromHtmlInfo(info, oembed), siteName)
	excerpt := cleanExcerpt(getExcerptFromHtmlInfo(info, oembed), title)
	fullText := cleanFullText(info.MainContent, title, excerpt)

	// { "__source__": { "html_info": ..., "oembed": ... }, title: "...", .... }
	data.Title = title
	data.Excerpt = excerpt
	data.FullText = fullText
	data.FullTextHash = generateHash(fullText)
	data.FoundAt = time.Now()
	data.Url = url
	data.UrlHash = generateHash(url)
	if info.OGInfo != nil && info.OGInfo.Article != nil {
		data.PublishedAt = info.OGInfo.Article.PublishedTime
		data.ModifiedAt = info.OGInfo.Article.ModifiedTime
	}
	data.Images = getOGImages(info.OGInfo.Images)
	data.ImageUrl = getImageUrl(oembed.ThumbnailURL, data.Images)

	// url = info.CanonicalUrl || url?
	// images = info.Images
	// image_url = oembed.ThumbnailUrl || info.Images.First that does not have cover in the name of img, otherwise just return cover
	// modified_at =
	// published_at =
	// updated_at =
	// TODO: Ability to override some stuff, like for example the published at
	// and all images from http://www.ba.gov.br/noticias/bahia-alcanca-segundo-lugar-no-ranking-nacional-de-testagens

	// log.Warn("Truncating full content!!!")
	// info.MainContent = info.MainContent[0:100]
	// data.FullText = fullText[0:200]

	return data, nil

	// TODO: * Remove <html>, <head>, <body>, trim elements, etc
	//       * Remove site_name from title
	//       * Remove excerpt and title from content
	//       * Give preference to images that don't have `cover` in the name
}

func generateHash(text string) string {
	algorithm := sha1.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func cleanFullText(html, title, excerpt string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		panic(err)
	}

	fullText := doc.Find("body").Text()
	fullTextChunks := []string{}
	for _, str := range strings.Split(fullText, "\n") {
		str = strings.Trim(str, " ")
		str = strings.Trim(str, "\t")
		if str != "" && str != title && str != excerpt {
			fullTextChunks = append(fullTextChunks, str)
		}
	}
	fullText = strings.Join(fullTextChunks, "\n")
	return fullText
}

func getSiteName(info *htmlinfo.HTMLInfo, oembed *oembed.Info) string {
	if info.OGInfo != nil && info.OGInfo.SiteName != "" {
		return info.OGInfo.SiteName
	}
	if oembed != nil && oembed.ProviderName != "" {
		return oembed.ProviderName
	}
	return ""
}

func getTitleFromHtmlInfo(info *htmlinfo.HTMLInfo, oembed *oembed.Info) string {
	if oembed.Title != "" {
		return oembed.Title
	}
	if info.Title != "" {
		return info.Title
	}

	return ""
}

func getExcerptFromHtmlInfo(info *htmlinfo.HTMLInfo, oembed *oembed.Info) string {
	if oembed.Description != "" {
		return oembed.Description
	}
	if info.Description != "" {
		return info.Description
	}

	return ""
}

func getOGImages(images []*opengraph.Image) []*core.ScrapedArticleImage {
	imgs := make([]*core.ScrapedArticleImage, len(images))
	for i, img := range images {
		imgs[i] = &core.ScrapedArticleImage{
			Url:       img.URL,
			SecureUrl: img.SecureURL,
			Type:      img.Type,
			Width:     img.Width,
			Height:    img.Height,
		}
	}
	return imgs
}

func getImageUrl(oembedThumb string, imgs []*core.ScrapedArticleImage) string {
	for _, img := range imgs {
		if img.Width != 0 && img.Height != 0 {
			if img.SecureUrl != "" {
				return img.SecureUrl
			}
			if img.Url != "" {
				return img.Url
			}
		}
	}
	return oembedThumb
}

func cleanTitle(title, siteName string) string {
	if title == "" || siteName == "" {
		return ""
	}

	// Use a regex here if this list grows, if none of this fixes the problem with
	// titles, use a generic regex
	title = strings.TrimSuffix(title, " - "+siteName)
	title = strings.TrimSuffix(title, " | "+siteName)
	// TODO: Remove source name too (Like for http://www.ba.gov.br/noticias/bahia-alcanca-segundo-lugar-no-ranking-nacional-de-testagens)

	return title
}

func cleanExcerpt(excerpt, title string) string {
	// TODO: Remove timestamps, like the ones in
	// http://www.amazonas.am.gov.br/2020/05/casa-do-migrante-jacamim-27-anos-acolhendo-pessoas-em-situacao-de-vulnerabilidade/
	return strings.Trim(strings.TrimPrefix(excerpt, title), " ")
}
