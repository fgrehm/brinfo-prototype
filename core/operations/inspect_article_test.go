package operations_test

import (
	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/operations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InspectArticle", func() {
	var (
		fakeData *ScrapedArticleData
		scraper  *fakeScraper
	)

	BeforeEach(func() {
		ts = newTestServer()
	})

	AfterEach(func() {
		ts.Close()
	})

	Context("validations", func() {
		It("fails if no url provided", func() {
			_, err := InspectArticle(InspectArticleInput{Url: ""})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^No URL provided$"))
		})
	})

	Context("extraction", func() {
		BeforeEach(func() {
			fakeData = &ScrapedArticleData{Title: "test title"}
			scraper = &fakeScraper{fakeData}
		})

		It("gets delegated to the article scraper", func() {
			data, err := InspectArticle(InspectArticleInput{
				Url:            ts.URL + "/good",
				ArticleScraper: scraper,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(data.ScrapedArticleData).To(Equal(fakeData))
		})

		It("returns an error if http response is not 200", func() {
			_, err := InspectArticle(InspectArticleInput{
				Url:            ts.URL + "/bad",
				ArticleScraper: scraper,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^Not Found$"))
		})
	})
})
