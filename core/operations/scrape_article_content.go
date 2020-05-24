package operations

import (
	"errors"
	"fmt"
	neturl "net/url"

	. "github.com/fgrehm/brinfo/core"
)

type ScrapeArticleContentInput struct {
	Repo          ContentSourceRepo
	ContentSource *ContentSource
	Url           string
}

func ScrapeArticleContent(input ScrapeArticleContentInput) (*ScrapedArticleData, error) {
	var (
		url  = input.Url
		cs   = input.ContentSource
		repo = input.Repo

		err error
	)

	if err = validateScrapeArticleContentInput(input); err != nil {
		return nil, err
	}

	if cs == nil {
		cs, err = lookupContentSourceForUrl(repo, url)
		if err != nil {
			return nil, err
		}
	}

	if err = validateContentSourceForScraping(cs, url); err != nil {
		return nil, err
	}

	body, contentType, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	data, err := cs.ArticleScraper.Run(body, url, contentType)
	if err != nil {
		return nil, err
	}
	data.SourceID = cs.ID
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

func lookupContentSourceForUrl(repo ContentSourceRepo, url string) (*ContentSource, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	return repo.FindByHost(u.Host)
}
