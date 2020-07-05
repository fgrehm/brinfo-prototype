package operations

import (
	"context"
	neturl "net/url"
	"time"

	// . "github.com/fgrehm/brinfo/core"

	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
)

func loggerFromContext(ctx context.Context) log.Interface {
	return log.FromContext(ctx)
}

type realClock struct{}

func (*realClock) Now() time.Time {
	return time.Now()
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

func extractHost(url string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}
