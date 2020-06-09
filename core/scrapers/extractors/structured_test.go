package extractors_test

import (
	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Structured", func() {
	It("works", func() {
		extractor := Structured("head", map[string]Extractor{
			"title":       Text("title", false),
			"description": Attribute("meta[name=description]", "content"),
		})

		val, err := extractor.Extract(Fragment(`
		<html>
			<head>
				<title>page title</title>
				<meta name="description" content="some description">
			</head>
			<body>
				foo
				<title>AAA</title>
				<meta name="description" content="FFF">
			</body>
		</html>`))

		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal(map[string]ExtractorResult{
			"title":       "page title",
			"description": "some description",
		}))

		_, err = extractor.Extract(Fragment(`<bad example=true>`))
		Expect(err).To(HaveOccurred())

		_, err = extractor.Extract(Fragment(`<head><title>AA</title></head>`))
		Expect(err).To(HaveOccurred())
	})

	It("works as a wrapper for multiple elements", func() {
		extractor := StructuredList("li", map[string]Extractor{
			"title":       Text("h2", false),
			"description": Text("p", false),
		})

		val, err := extractor.Extract(Fragment(`
		<ul>
			<li>
				<h2>Title 1</h2>
				<p>Desc 1</p>
				<em>bla</em>
			</li>
			<li>
				<h2>Title 2</h2>
				<p>Desc 2</p>
				<em>bla 2</em>
			</li>
		</ul>`))

		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal([]map[string]ExtractorResult{
			{"title": "Title 1", "description": "Desc 1"},
			{"title": "Title 2", "description": "Desc 2"},
		}))

		_, err = extractor.Extract(Fragment(`<bad example=true>`))
		Expect(err).To(HaveOccurred())

		_, err = extractor.Extract(Fragment(`<head><title>AA</title></head>`))
		Expect(err).To(HaveOccurred())
	})
})
