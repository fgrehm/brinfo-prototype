package core_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/fgrehm/brinfo/core"
)

var _ = Describe("Core", func() {
	Context("ScrapedArticleData", func() {
		Context("ValidForIngestion", func() {
			var data *ScrapedArticleData

			BeforeEach(func() {
				now := time.Now()
				data = &ScrapedArticleData{
					SourceID:     "some-source",
					ContentType:  "article",
					Url:          "https://example.com",
					UrlHash:      "url-hash",
					Title:        "Article title",
					FullText:     "text",
					FullTextHash: "text-hash",
					FoundAt:      now,
					PublishedAt:  &now,
					ModifiedAt:   &now,
					Images: []*ScrapedArticleImage{
						&ScrapedArticleImage{Url: "http://image.url"},
					},
					ImageUrl: "http://image.url",
				}

				Expect(data.ValidForIngestion()).To(BeTrue())
			})

			It("is invalid if no SourceID is set", func() {
				data.SourceID = ""
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is invalid if no ContentType is set", func() {
				data.ContentType = ""
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is invalid if no Url is set", func() {
				data.Url = ""
				Expect(data.ValidForIngestion()).To(BeFalse())
			})

			It("is invalid if no UrlHash is set", func() {
				data.UrlHash = ""
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

			It("is valid if no Images are set", func() {
				data.Images = nil
				Expect(data.ValidForIngestion()).To(BeTrue())
			})

			It("is valid if no ImageUrl is set", func() {
				data.ImageUrl = ""
				Expect(data.ValidForIngestion()).To(BeTrue())
			})
		})
	})
})
