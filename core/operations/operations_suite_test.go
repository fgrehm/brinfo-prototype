package operations_test

import (
	"io"
	"net/http"
	"testing"

	. "github.com/fgrehm/brinfo/core"
	op "github.com/fgrehm/brinfo/core/operations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http/httptest"
)

var (
	ts *httptest.Server
)

func TestOperations(t *testing.T) {
	op.UseCache = false
	RegisterFailHandler(Fail)
	RunSpecs(t, "Operations Suite")
}

func newTestServer() *httptest.Server {
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

	return httptest.NewServer(mux)
}

type fakeScraper struct {
	data *ScrapedArticleData
}

func (f *fakeScraper) Run(io.Reader, string) (*ScrapedArticleData, error) {
	return f.data, nil
}
