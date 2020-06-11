package operations

import (
	"context"
	"fmt"
	neturl "net/url"
	"time"

	. "github.com/fgrehm/brinfo/core"

	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
)

func loggerFromContext(ctx context.Context) log.Interface {
	return log.FromContext(ctx)
}

func doScrapeArticle(ctx context.Context, cache bool, url string, cs *ContentSource, scraper ArticleScraper) (*ScrapedArticleData, error) {
	log := loggerFromContext(ctx)

	body, contentType, err := makeRequest(cache, url)
	if err != nil {
		return nil, err
	}

	if cs != nil && cs.ForceContentType != "" {
		log.Debugf("Forcing content type to %s", cs.ForceContentType)
		contentType = cs.ForceContentType
	}

	data, err := InspectBytes(ctx, InspectBytesArgs{
		HTML:           body,
		URL:            url,
		ContentType:    contentType,
		ArticleScraper: scraper,
	})
	if cs != nil {
		data.SourceID = cs.ID
	}
	return data, nil
}

func validateContentSourceForScraping(cs *ContentSource, url string) error {
	if cs.Host == "" {
		return fmt.Errorf("ContentSource does not have a host set %+v", cs)
	}
	if cs.ArticleScraper == nil {
		return fmt.Errorf("Article scraper not assigned for ContentSource '%+v'", cs)
	}

	host, err := extractHost(url)
	if err != nil {
		return err
	}
	if host != cs.Host {
		return fmt.Errorf("URL host '%s' does not match ContentSource host '%s'", host, cs.Host)
	}

	return nil
}

func makeRequest(cache bool, url string) ([]byte, string, error) {
	opts := []colly.CollectorOption{
		colly.UserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:76.0) Gecko/20100101 Firefox/76.0"),
	}

	if cache {
		log.Info("Using cache")
		opts = append(opts, colly.CacheDir("./.brinfo-cache/"))
	}
	c := colly.NewCollector(opts...)
	c.SetRequestTimeout(5 * time.Second)

	var (
		body        []byte
		contentType string
		err         error
	)

	c.OnResponse(func(r *colly.Response) {
		log.Debugf("Status: %d", r.StatusCode)
		if r.StatusCode == 200 {
			body = r.Body
			contentType = r.Headers.Get("Content-Type")
		}
	})

	err = c.Visit(url)
	if err != nil {
		return nil, "", err
	}
	c.Wait()

	if err != nil {
		return nil, "", err
	}

	return body, contentType, nil
}

func mustLookupContentSourceForUrl(ctx context.Context, repo ContentSourceRepo, url string) (*ContentSource, error) {
	host, err := extractHost(url)
	if err != nil {
		panic(err)
	}
	return repo.FindByHost(ctx, host)
}

func lookupContentSourceForUrl(ctx context.Context, repo ContentSourceRepo, url string) (*ContentSource, error) {
	host, err := extractHost(url)
	if err != nil {
		return nil, err
	}
	return repo.GetByHost(ctx, host)
}

func extractHost(url string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}
