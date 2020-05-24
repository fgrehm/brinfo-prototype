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
})
