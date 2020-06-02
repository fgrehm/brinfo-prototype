package core_test

import (
	"context"
	"errors"
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

	Context("CombinedArticleScraper", func() {
		var (
			yesterday time.Time
			now       time.Time
			baseData  *ScrapedArticleData
			moreData  *ScrapedArticleData
		)

		BeforeEach(func() {
			now = time.Now()
			yesterday = now.AddDate(0, 0, -1)

			baseData = &ScrapedArticleData{
				Extra:        map[string]interface{}{"a": 1},
				SourceID:     "fid",
				ContentType:  "",
				Url:          "f-url",
				UrlHash:      "f-url-hash",
				Title:        "f-title",
				FullText:     "f text",
				FullTextHash: "f-text-hash",
				Excerpt:      "f excerpt",
				FoundAt:      yesterday,
				PublishedAt:  &yesterday,
				ModifiedAt:   &yesterday,
				Images: []*ScrapedArticleImage{
					&ScrapedArticleImage{Url: "img-url"},
				},
				ImageUrl: "last-img-url",
			}

			moreData = &ScrapedArticleData{
				Extra:        map[string]interface{}{"b": 2},
				SourceID:     "lastid",
				ContentType:  "article",
				Url:          "last-url",
				UrlHash:      "last-url-hash",
				Title:        "last-title",
				FullText:     "last text",
				FullTextHash: "last-text-hash",
				Excerpt:      "last excerpt",
				FoundAt:      now,
				PublishedAt:  &now,
				ModifiedAt:   &now,
				Images: []*ScrapedArticleImage{
					&ScrapedArticleImage{Url: "img-url2"},
				},
				ImageUrl: "last-img-url",
			}
		})

		It("returns the data from all scrapers merged", func() {
			firstScraper := &fakeScraper{result: baseData}
			lastScraper := &fakeScraper{result: moreData}

			scraper := CombinedArticleScraper(firstScraper, lastScraper)

			data, err := scraper.Run(context.Background(), []byte{}, "", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(data).To(Equal(&ScrapedArticleData{
				Extra: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
				SourceID:     "lastid",
				ContentType:  "article",
				Url:          "last-url",
				UrlHash:      "last-url-hash",
				Title:        "last-title",
				FullText:     "last text",
				FullTextHash: "last-text-hash",
				Excerpt:      "last excerpt",
				FoundAt:      now,
				PublishedAt:  &now,
				ModifiedAt:   &now,
				Images: []*ScrapedArticleImage{
					&ScrapedArticleImage{Url: "img-url"},
					&ScrapedArticleImage{Url: "img-url2"},
				},
				ImageUrl: "last-img-url",
			}))
		})

		It("last provided scraper data wins, if value set", func() {
			baseData.ContentType = "article"
			firstScraper := &fakeScraper{result: baseData}
			moreData.ContentType = ""
			bestScraper := &fakeScraper{result: moreData}

			scraper := CombinedArticleScraper(firstScraper, &fakeScraper{&ScrapedArticleData{}, nil}, bestScraper, &fakeScraper{&ScrapedArticleData{}, nil})

			data, err := scraper.Run(context.Background(), []byte{}, "", "")
			Expect(err).NotTo(HaveOccurred())
			Expect(data).To(Equal(&ScrapedArticleData{
				Extra: map[string]interface{}{
					"a": 1,
					"b": 2,
				},
				SourceID:     "lastid",
				ContentType:  "article",
				Url:          "last-url",
				UrlHash:      "last-url-hash",
				Title:        "last-title",
				FullText:     "last text",
				FullTextHash: "last-text-hash",
				Excerpt:      "last excerpt",
				FoundAt:      now,
				PublishedAt:  &now,
				ModifiedAt:   &now,
				Images: []*ScrapedArticleImage{
					&ScrapedArticleImage{Url: "img-url"},
					&ScrapedArticleImage{Url: "img-url2"},
				},
				ImageUrl: "last-img-url",
			}))
		})

		It("fails if any of the scraper fails", func() {
			firstScraper := &fakeScraper{result: baseData}
			bestScraper := &fakeScraper{result: moreData}
			errScraper := &fakeScraper{err: errors.New("BOOM")}

			scraper := CombinedArticleScraper(firstScraper, errScraper, bestScraper, &fakeScraper{&ScrapedArticleData{}, nil})

			data, err := scraper.Run(context.Background(), []byte{}, "", "")
			Expect(err).To(HaveOccurred())
			Expect(data).To(BeNil())
		})
	})
})

type fakeScraper struct {
	result *ScrapedArticleData
	err    error
}

func (s *fakeScraper) Run(context.Context, []byte, string, string) (*ScrapedArticleData, error) {
	return s.result, s.err
}
