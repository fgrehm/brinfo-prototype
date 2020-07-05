package operations_test

import (
	"context"
	"time"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/operations"

	"github.com/fgrehm/brinfo/core/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ScrapeArticlesListing", func() {
	var (
		ctx context.Context
		ts  *testutils.Server
	)

	brLoc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	BeforeEach(func() {
		ctx = context.Background()
		ts = testutils.NewTestServer()
	})

	AfterEach(func() {
		ts.Close()
	})

	It("fetches URLs based on a selector", func() {
		ts.Articles = []*testutils.Article{
			{URL: ts.URL() + "/first-article"},
			{URL: ts.URL() + "/second-article"},
		}

		result, err := ScrapeArticlesListing(ctx, ScrapeArticlesListingArgs{
			URL:           ts.URL() + "/articles",
			LinkContainer: "ul li",
			URLExtractor:  "a[href] | href",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal([]*ArticleLink{
			{URL: ts.URL() + "/first-article"},
			{URL: ts.URL() + "/second-article"},
		}))
	})

	It("fixes relative URLs", func() {
		ts.Articles = []*testutils.Article{
			{URL: "/articles/first-article"},
			{URL: "second-article"},
		}

		result, err := ScrapeArticlesListing(ctx, ScrapeArticlesListingArgs{
			URL:           ts.URL() + "/articles",
			LinkContainer: "ul li",
			URLExtractor:  "a[href] | href",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal([]*ArticleLink{
			{URL: ts.URL() + "/articles/first-article"},
			{URL: ts.URL() + "/second-article"},
		}))
	})

	It("supports extraction of article metadata", func() {
		ts.Articles = []*testutils.Article{
			{URL: "first-article", ImageURL: "/img.png"},
			{URL: "second-article", PublishedAt: "08/06/2020 23:11"},
		}

		sampleDate := time.Date(2020, 6, 8, 23, 11, 0, 0, brLoc)
		sampleImg := ts.URL() + ts.Articles[0].ImageURL

		result, err := ScrapeArticlesListing(ctx, ScrapeArticlesListingArgs{
			URL:                  ts.URL() + "/articles",
			LinkContainer:        "ul li",
			URLExtractor:         "a[href] | href",
			PublishedAtExtractor: "time | text?::time",
			ImageURLExtractor:    "img | src?",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal([]*ArticleLink{
			{URL: ts.URL() + "/first-article", PublishedAt: nil, ImageURL: &sampleImg},
			{URL: ts.URL() + "/second-article", PublishedAt: &sampleDate, ImageURL: nil},
		}))
	})
})
