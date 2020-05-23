package operations

import (
	"bytes"
	"io/ioutil"
	neturl "net/url"
	"time"

	"github.com/dimchansky/utfbom"
	"github.com/gocolly/colly/v2"
)

var UseCache = true

func makeRequest(url string) ([]byte, error) {
	opts := []colly.CollectorOption{
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	}

	if UseCache {
		opts = append(opts, colly.CacheDir("./.brinfo-cache/"))
	}
	c := colly.NewCollector(opts...)
	c.SetRequestTimeout(5 * time.Second)

	var (
		body []byte
		err  error
	)

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode == 200 {
			body = r.Body
		}
	})

	err = c.Visit(url)
	if err != nil {
		return nil, err
	}
	c.Wait()

	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(utfbom.SkipOnly(bytes.NewBuffer(body)))
	if err != nil {
		return nil, err
	}

	return body, nil
}

func extractHost(url string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}
