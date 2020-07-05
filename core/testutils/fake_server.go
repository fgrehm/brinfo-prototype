package testutils

import (
	"net/http"
	"net/http/httptest"
)

type Server struct {
	server   *httptest.Server
	Articles []*Article
	PerPage  int
}

type Article struct {
	ID          string
	URL         string
	ImageURL    string
	Title       string
	PublishedAt string
	Excerpt     string
	Head        string
	Body        string
}

func NewTestServer() *Server {
	ts := &Server{PerPage: 5, Articles: []*Article{}}
	mux := http.NewServeMux()

	mux.HandleFunc("/articles", ts.listArticles)
	mux.HandleFunc("/articles/show", ts.showArticle)

	ts.server = httptest.NewServer(mux)
	return ts
}

func (s *Server) URL() string {
	return s.server.URL
}

func (s *Server) Close() {
	s.server.Close()
}

func (s *Server) listArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
<html>
	<head>
		<title>Articles</title>
	</head>
	<body>
		<h1>All articles</h1>
		<ul>
			` + s.renderArticlesList() + `
		</ul>
	</body>
</html>`))
}

func (s *Server) renderArticlesList() string {
	list := ""
	for _, a := range s.Articles {
		linkText := a.Title
		if linkText == "" {
			linkText = a.URL
		}
		list += `<li><a href="` + a.URL + `">` + linkText + `</a>`
		if a.PublishedAt != "" {
			list += `<time>` + a.PublishedAt + `</time>`
		}
		if a.ImageURL != "" {
			list += `<img src="` + a.ImageURL + `">`
		}
		list += `</li>`
	}

	return list
}

func (s *Server) showArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	panic("BOOOM")
}

// func getArticleByID
