package core_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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

				Expect(data.ValidForIngestion()).To(BeTrue())
			})

			It("is invalid if no URL is set", func() {
				data.URL = ""
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is invalid if no URLHash is set", func() {
				data.URLHash = ""
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is invalid if no Title is set", func() {
				data.Title = ""
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is invalid if no FullText is set", func() {
				data.FullText = ""
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is invalid if no FullTextHash is set", func() {
				data.FullTextHash = ""
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is invalid if no FoundAt is set", func() {
				var def time.Time
				data.FoundAt = def
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is invalid if no PublishedAt is set", func() {
				data.PublishedAt = nil
				Expect(data.ValidForIngestion()).To(BeFalse())

				var def time.Time
				data.PublishedAt = &def
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is valid if no Excerpt is set", func() {
				data.Excerpt = ""
				Expect(data.ValidForIngestion()).To(BeTrue())
			})

			It("is invalid if ModifiedAt is set to default", func() {
				var def time.Time
				data.ModifiedAt = &def
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is valid if no ModifiedAt is set", func() {
				data.ModifiedAt = nil
				Expect(data.ValidForIngestion()).To(BeTrue())
			})

			It("is valid if no ImageURL is set", func() {
				data.ImageURL = ""
				Expect(data.ValidForIngestion()).To(BeTrue())
			})
		})
	})
})
