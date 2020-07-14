package core_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gomegat "github.com/onsi/gomega/types"

	. "github.com/fgrehm/brinfo/core"
)

var _ = Describe("Core", func() {
	Context("ArticleData", func() {
		Context("ValidForIngestion", func() {
			var data *ArticleData

			BeforeEach(func() {
				now := time.Now()
				data = &ArticleData{
					URL:          "https://example.com",
					URLHash:      "url-hash",
					Title:        "Article title",
					FullText:     "text",
					FullTextHash: "text-hash",
					FoundAt:      now,
					PublishedAt:  &now,
					ModifiedAt:   &now,
					ImageURL:     "http://image.url",
				}

				Expect(data).To(BeValidForIngestion())
			})

			It("is invalid if no URL is set", func() {
				data.URL = ""
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is invalid if no URLHash is set", func() {
				data.URLHash = ""
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is invalid if no Title is set", func() {
				data.Title = ""
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is invalid if no FullText is set", func() {
				data.FullText = ""
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is invalid if no FullTextHash is set", func() {
				data.FullTextHash = ""
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is invalid if no FoundAt is set", func() {
				var def time.Time
				data.FoundAt = def
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is invalid if no PublishedAt is set", func() {
				data.PublishedAt = nil
				Expect(data).NotTo(BeValidForIngestion())

				var def time.Time
				data.PublishedAt = &def
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is invalid if PublishedAt is in the future", func() {
				newTime := data.PublishedAt.Add(time.Hour * 24)
				data.PublishedAt = &newTime
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is valid if no Excerpt is set", func() {
				data.Excerpt = ""
				Expect(data).To(BeValidForIngestion())
			})

			It("is invalid if ModifiedAt is set to default", func() {
				var def time.Time
				data.ModifiedAt = &def
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is invalid if ModifiedAt is in the future", func() {
				newTime := data.PublishedAt.Add(time.Hour * 24)
				data.ModifiedAt = &newTime
				Expect(data).NotTo(BeValidForIngestion())
			})

			It("is valid if no ModifiedAt is set", func() {
				data.ModifiedAt = nil
				Expect(data).To(BeValidForIngestion())
			})

			It("is valid if no ImageURL is set", func() {
				data.ImageURL = ""
				Expect(data).To(BeValidForIngestion())
			})

			It("is valid if no ImageURL is relative", func() {
				data.ImageURL = "/foo/bar"
				Expect(data).NotTo(BeValidForIngestion())

				data.ImageURL = "fooo.com/foo/bar"
				Expect(data).NotTo(BeValidForIngestion())
			})
		})
	})
})

type validForIngestionMatcher struct{}

func BeValidForIngestion() gomegat.GomegaMatcher {
	return &validForIngestionMatcher{}
}

func (m *validForIngestionMatcher) Match(actual interface{}) (bool, error) {
	data, ok := actual.(*ArticleData)
	if !ok {
		return false, fmt.Errorf("BeValidForIngestion matcher expects a pointer to core.ArticleData")
	}

	valid, _ := data.ValidForIngestion()
	return valid, nil
}

func (m *validForIngestionMatcher) FailureMessage(actual interface{}) string {
	data, _ := actual.(*ArticleData)
	_, errMsgs := data.ValidForIngestion()
	return fmt.Sprintf("Expected data to be valid for ingestion but it is not: %v", errMsgs)
}

func (m *validForIngestionMatcher) NegatedFailureMessage(actual interface{}) string {
	return "Expected data to be invalid for ingestion"
}
