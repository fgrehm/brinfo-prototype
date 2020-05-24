package operations

import (
	neturl "net/url"
	"time"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

var UseCache = true

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

func extractHost(url string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}
