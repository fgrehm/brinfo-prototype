package operations_test

import (
	"context"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/operations"
	mem "github.com/fgrehm/brinfo/storage/inmemory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InspectArticle", func() {
	var (
		fakeData *ScrapedArticleData
		scraper  *fakeScraper
		ctx      context.Context
	)

	BeforeEach(func() {
		ts = newTestServer()
		ctx = context.Background()
	})

	AfterEach(func() {
		ts.Close()
	})

	Context("validations", func() {
		It("fails if no url provided", func() {
			_, err := InspectArticle(ctx, InspectArticleInput{Url: "", ContentSourceRepo: mem.NewContentSourceRepo()})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^No URL provided$"))
		})

		It("fails if no content repo provided", func() {
			_, err := InspectArticle(ctx, InspectArticleInput{Url: "http://go.com"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^No content source repo provided$"))
		})
	})

	Context("extraction", func() {
		BeforeEach(func() {
			fakeData = &ScrapedArticleData{Title: "test title"}
			scraper = &fakeScraper{fakeData}
		})

		It("gets delegated to the article scraper", func() {
			data, err := InspectArticle(ctx, InspectArticleInput{
				Url:               ts.URL + "/good",
				ArticleScraper:    scraper,
				ContentSourceRepo: mem.NewContentSourceRepo(),
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(data.ScrapedArticleData).To(Equal(fakeData))
		})

		It("returns an error if http response is not 200", func() {
			_, err := InspectArticle(ctx, InspectArticleInput{
				Url:               ts.URL + "/bad",
				ArticleScraper:    scraper,
				ContentSourceRepo: mem.NewContentSourceRepo(),
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^Not Found$"))
		})
	})
})
