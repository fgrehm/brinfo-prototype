package extractors_test

import (
	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Attribute", func() {
	It("works", func() {
		e := Attribute("meta", "name")

		val, err := e.Extract(Fragment(`<meta name="bla">`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("bla"))

		_, err = e.Extract(Fragment(`<metas name="bla">`))
		Expect(err).To(HaveOccurred())

		_, err = e.Extract(Fragment(`<meta names="bla">`))
		Expect(err).To(HaveOccurred())

		_, err = e.Extract(Fragment(`<meta name="bla"><meta name="bla2">`))
		Expect(err).To(HaveOccurred())
	})
})

// 	It("can extract multiple attributes with ExtractMultipleAttributes", func() {
// 		scraper := NewScraper().ExtractMultipleAttributes("meta", "meta[name*=og]", "content")

// 		result, err := scraper.Run(`<meta name="og:a" content="a"><meta name="og:b" content="b">`)
// 		Expect(err).NotTo(HaveOccurred())

// 		val, err := result.Get("meta")
// 		Expect(err).NotTo(HaveOccurred())
// 		Expect(val).To(Equal([]string{"a", "b"}))
// 	})

// 	XIt("can be made optional with ExtractOptionalAttribute", func() {
// 	})
