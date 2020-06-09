package operations_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/fgrehm/brinfo/core"
	op "github.com/fgrehm/brinfo/core/operations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOperations(t *testing.T) {
	op.UseCache = false
	RegisterFailHandler(Fail)
	RunSpecs(t, "Operations Suite")
}

type testServer struct {
	server   *httptest.Server
	articles []*testArticle
	perPage  int
}

type testArticle struct {
	url         string
	imageURL    string
	title       string
	publishedAt string
	excerpt     string
	body        string
}

func newTestServer() *testServer {
	ts := &testServer{perPage: 5}
	mux := http.NewServeMux()

	mux.HandleFunc("/articles", ts.listArticles)
	mux.HandleFunc("/good", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
	<head>
		<title>Test Page</title>
	</head>
	<body>
		<h1>Hello World</h1>
		<p class="description">This is a test page</p>
		<p class="description">This is a test paragraph</p>
	</body>
</html>
		`))
	})

	ts.server = httptest.NewServer(mux)
	return ts
}

func (s *testServer) URL() string {
	return s.server.URL
}

func (s *testServer) Close() {
	s.server.Close()
}

func (s *testServer) listArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html>
	<head>
		<title>Articles</title>
	</head>
	<body>
		<h1>Page header</h1>
		<ul>
			` + s.renderArticlesList() + `
		</ul>
	</body>
</html>`))
}

func (s *testServer) renderArticlesList() string {
	if s.articles == nil || len(s.articles) == 0 {
		// TODO: Set a default array from newTestServer()
		return `
<li>
	<a href="http://example.com">Out same tab</a>
</li>
<li>
	<a href="http://example.com" target="_blank">Out new tab</a>
</li>
<li>
	<a href="/article?id=1">Relative</a>
</li>
<li>
	<a href="` + s.URL() + `/article?id=2">Full path</a>
</li>`
	}

	list := ""
	for _, a := range s.articles {
		linkText := a.title
		if linkText == "" {
			linkText = a.url
		}
		list += `<li><a href="` + a.url + `">` + linkText + `</a>`
		if a.publishedAt != "" {
			list += `<time>` + a.publishedAt + `</time>`
		}
		if a.imageURL != "" {
			list += `<img src="` + a.imageURL + `">`
		}
		list += `</li>`
	}

	return list
}

type fakeScraper struct {
	data *ScrapedArticleData
}

func (f *fakeScraper) Run(context.Context, []byte, string, string) (*ScrapedArticleData, error) {
	return f.data, nil
}
