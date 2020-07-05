package extractors_test

import (
	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Structured", func() {
	It("works", func() {
		e := Structured("head", map[string]Extractor{
			"title":       Text("title", false),
			"description": Attribute("meta[name=description]", "content"),
		})

		val, err := extract(e, `<html>
			<head>
				<title>page title</title>
				<meta name="description" content="some description">
			</head>
			<body>
				foo
				<title>AAA</title>
				<meta name="description" content="FFF">
			</body>
		</html>`)

		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal(map[string]ExtractorResult{
			"title":       "page title",
			"description": "some description",
		}))

		_, err = extract(e, `<bad example=true>`)
		Expect(err).To(HaveOccurred())

		_, err = extract(e, `<head><title>AA</title></head>`)
		Expect(err).To(HaveOccurred())
	})

	It("works as a wrapper for multiple elements", func() {
		e := StructuredList("li", map[string]Extractor{
			"title":       Text("h2", false),
			"description": Text("p", false),
		})

		val, err := extract(e, `
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
		</ul>`)

		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal([]map[string]ExtractorResult{
			{"title": "Title 1", "description": "Desc 1"},
			{"title": "Title 2", "description": "Desc 2"},
		}))

		_, err = extract(e, `<bad example=true>`)
		Expect(err).To(HaveOccurred())

		_, err = extract(e, `<head><title>AA</title></head>`)
		Expect(err).To(HaveOccurred())
	})
})
