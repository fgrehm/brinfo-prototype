package extractors_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExtractors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Extractors Suite")
}

func Fragment(html string) *goquery.Selection {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	Expect(err).NotTo(HaveOccurred())
	return doc.Selection
}
