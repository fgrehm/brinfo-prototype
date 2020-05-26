package operations

import (
	"errors"
	"fmt"
	neturl "net/url"
	"time"

	. "github.com/fgrehm/brinfo/core"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

var UseCache = true

func doScrapeArticleContent(url string, cs *ContentSource, scraper ArticleScraper) (*ScrapedArticleData, error) {
	body, contentType, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	if cs != nil && cs.ForceContentType != "" {
		contentType = cs.ForceContentType
	}

	data, err := scraper.Run(body, url, contentType)
	if err != nil {
		return nil, err
	}
	if cs != nil {
		data.SourceID = cs.ID
	}
	return data, nil
}

func validateScrapeArticleContentInput(input ScrapeArticleContentInput) error {
	if input.Url == "" {
		return errors.New("No URL provided")
	}
	if input.ContentSource == nil && input.Repo == nil {
		return errors.New("No ContentSource or Repository provided")
	}

	return nil
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

func makeRequest(url string) ([]byte, string, error) {
	opts := []colly.CollectorOption{
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	}

	if UseCache {
		log.Debug("Using cache")
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

func mustLookupContentSourceForUrl(repo ContentSourceRepo, url string) (*ContentSource, error) {
	host, err := extractHost(url)
	if err != nil {
		return nil, err
	}
	return repo.FindByHost(host)
}

func lookupContentSourceForUrl(repo ContentSourceRepo, url string) (*ContentSource, error) {
	host, err := extractHost(url)
	if err != nil {
		return nil, err
	}
	return repo.GetByHost(host)
}

func extractHost(url string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}
