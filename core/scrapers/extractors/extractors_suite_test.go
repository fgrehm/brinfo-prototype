package extractors_test

import (
	"context"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExtractors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Extractors Suite")
}

// TODO: mustExtract to avoid checking for errors on each test

func extract(ext Extractor, html string) (ExtractorResult, error) {
	return extractURL(ext, "https://brinfo.io", html)
}

func extractURL(ext Extractor, url, html string) (ExtractorResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	return ext.Extract(ExtractorArgs{
		Context: context.Background(),
		URL:     url,
		Root:    doc.Selection,
	})
}
