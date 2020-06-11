package operations_test

import (
	"context"
	neturl "net/url"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/operations"
	mem "github.com/fgrehm/brinfo/storage/inmemory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ScrapeArticle", func() {
	var (
		repo     ContentSourceRepo
		cs       *ContentSource
		fakeData *ScrapedArticleData
		ctx      context.Context
		ts       *testServer
	)

	BeforeEach(func() {
		ctx = context.Background()
		ts = newTestServer()

		u, err := neturl.Parse(ts.URL())
		Expect(err).NotTo(HaveOccurred())

		fakeData = &ScrapedArticleData{Title: "test title"}
		cs = &ContentSource{
			ID:   "test-source",
			Host: u.Host,
		}

		repo = mem.NewContentSourceRepo()
		err = repo.Register(cs)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		ts.Close()
	})

	Context("validations", func() {
		It("fails if no url provided", func() {
			_, err := ScrapeArticle(ctx, ScrapeArticleArgs{
				URL:  "",
				Repo: mem.NewContentSourceRepo(),
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^No URL provided$"))

			_, err = ScrapeArticle(ctx, ScrapeArticleArgs{
				URL:           "",
				ContentSource: cs,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^No URL provided$"))
		})

		It("fails if no content source and repo are provided", func() {
			_, err := ScrapeArticle(ctx, ScrapeArticleArgs{URL: ts.URL()})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^No ContentSource or Repository provided$"))
		})

		It("fails if content source does not have an article scraper", func() {
			_, err := ScrapeArticle(ctx, ScrapeArticleArgs{
				URL:           ts.URL(),
				ContentSource: cs,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^Article scraper not assigned for ContentSource"))
		})

		It("fails if URL does not match content source host", func() {
			_, err := ScrapeArticle(ctx, ScrapeArticleArgs{
				URL: ts.URL(),
				ContentSource: &ContentSource{
					Host:           "http://example.com",
					ArticleScraper: &fakeScraper{},
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^URL host.*does not match ContentSource host.*$"))
		})
	})

	Context("extraction", func() {
		BeforeEach(func() {
			cs.ArticleScraper = &fakeScraper{fakeData}
		})

		It("gets delegated to ContentSource article scraper", func() {
			data, err := ScrapeArticle(ctx, ScrapeArticleArgs{
				URL:           ts.URL() + "/good",
				ContentSource: cs,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(data).To(Equal(fakeData))
		})

		It("assigns the content source ID", func() {
			Expect(fakeData.SourceID).To(BeEmpty())

			data, err := ScrapeArticle(ctx, ScrapeArticleArgs{
				URL:           ts.URL() + "/good",
				ContentSource: cs,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(data.SourceID).To(Equal(cs.ID))
		})

		It("returns an error if http response is not 200", func() {
			_, err := ScrapeArticle(ctx, ScrapeArticleArgs{
				URL:           ts.URL() + "/bad",
				ContentSource: cs,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(MatchRegexp("^Not Found$"))
		})

		Context("content source lookup", func() {
			It("works if repo knows the source for the given host", func() {
				data, err := ScrapeArticle(ctx, ScrapeArticleArgs{
					URL:  ts.URL() + "/good",
					Repo: repo,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(data).To(Equal(fakeData))
			})

			It("fails if repo doesnt know the source for the given host", func() {
				_, err := ScrapeArticle(ctx, ScrapeArticleArgs{
					URL:  "https://google.com",
					Repo: repo,
				})
				Expect(err).To(MatchError("Content source not found: google.com"))
			})
		})
	})
})
