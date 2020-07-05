package scrapers

import (
	"context"
	"time"

	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArticleScraper", func() {
	var (
		s   ArticleScraper
		cfg *ArticleScraperConfig
		ctx context.Context
		now time.Time
	)

	brLoc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	BeforeEach(func() {
		now = time.Now()
		cfg = &ArticleScraperConfig{Clock: fakeClock{now}}
		s = NewArticleScraper(cfg)
		ctx = context.Background()
	})

	It("sets the URL and FoundAt fields by default", func() {
		cfg.Extractors = []Extractor{&fakeExtractor{}}

		body := `<html><body><p>Don't care</p></body><html>`
		data, err := s.Run(ctx, []byte(body), "http://example.com", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(Equal(&ArticleData{
			Extra: map[string]interface{}{
				"html": mustGzip([]byte(body)),
			},
			URL:     "http://example.com",
			URLHash: "89dce6a446a69d6b9bdc01ac75251e4c322bcdff",
			FoundAt: now,
		}))
	})

	It("maps extractor result into an ArticleData", func() {
		cfg.Extractors = []Extractor{
			&fakeExtractor{map[string]interface{}{
				"title":   "Finally a cure for COVID19!",
				"excerpt": "A summary of how it attacks the virus",
			}},
		}

		body := `<html><body><p>Don't care</p></body><html>`
		data, err := s.Run(ctx, []byte(body), "http://example.com", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(Equal(&ArticleData{
			Extra: map[string]interface{}{
				"html": mustGzip([]byte(body)),
			},
			URL:     "http://example.com",
			URLHash: "89dce6a446a69d6b9bdc01ac75251e4c322bcdff",
			Title:   "Finally a cure for COVID19!",
			Excerpt: "A summary of how it attacks the virus",
			FoundAt: now,
		}))
	})

	It("combines data from multiple extractors into an ArticleData, last one wins", func() {
		cfg.Extractors = []Extractor{
			&fakeExtractor{map[string]interface{}{
				"title": "Will discard",
			}},
			&fakeExtractor{map[string]interface{}{
				"excerpt": "Random stuff",
			}},
			&fakeExtractor{map[string]interface{}{
				"title":   "Title",
				"excerpt": nil,
			}},
		}

		body := `<html><body><p>Don't care</p></body><html>`
		data, err := s.Run(ctx, []byte(body), "http://example.com", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(Equal(&ArticleData{
			Extra: map[string]interface{}{
				"html": mustGzip([]byte(body)),
			},
			URL:     "http://example.com",
			URLHash: "89dce6a446a69d6b9bdc01ac75251e4c322bcdff",
			Title:   "Title",
			Excerpt: "Random stuff",
			FoundAt: now,
		}))
	})

	It("is capable of extracting all of necessary ArticleData attributes", func() {
		pubDate := time.Date(2020, 6, 15, 19, 56, 0, 0, brLoc)
		modDate := time.Date(2020, 6, 15, 19, 59, 0, 0, brLoc)
		body := `<html><body><p>Don't care</p></body><html>`
		expectedData := &ArticleData{
			Extra: map[string]interface{}{
				"a":    "b",
				"html": mustGzip([]byte(body)),
			},
			URL:          "http://example.com",
			URLHash:      "89dce6a446a69d6b9bdc01ac75251e4c322bcdff",
			Title:        "Article title",
			FullText:     "Lots of text here",
			FullTextHash: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
			ImageURL:     "https://image.com",
			PublishedAt:  &pubDate,
			ModifiedAt:   &modDate,
			FoundAt:      now,
		}

		cfg.Extractors = []Extractor{
			&fakeExtractor{map[string]interface{}{
				"extra": map[string]interface{}{
					"a": "b",
				},
				"title":       expectedData.Title,
				"fullText":    expectedData.FullText,
				"imageURL":    expectedData.ImageURL,
				"publishedAt": pubDate,
				"modifiedAt":  modDate,
			}},
		}

		data, err := s.Run(ctx, []byte(body), "http://example.com", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(data).To(Equal(expectedData))
	})
})

type fakeClock struct {
	now time.Time
}

func (c fakeClock) Now() time.Time {
	return c.now
}

type fakeExtractor struct {
	result map[string]interface{}
}

func (e *fakeExtractor) Extract(ExtractorArgs) (ExtractorResult, error) {
	return e.result, nil
}
