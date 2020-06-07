package operations_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	slug        string
	title       string
	publishedAt time.Time
	excerpt     string
	body        string
}

func newTestServer() *testServer {
	mux := http.NewServeMux()

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

	return &testServer{
		server: httptest.NewServer(mux),
	}
}

func (s *testServer) URL() string {
	return s.server.URL
}

func (s *testServer) Close() {
	s.server.Close()
}

type fakeScraper struct {
	data *ScrapedArticleData
}

func (f *fakeScraper) Run(context.Context, []byte, string, string) (*ScrapedArticleData, error) {
	return f.data, nil
}
