package extractors_test

import (
	. "github.com/fgrehm/brinfo/core"
	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Images", func() {
	It("works", func() {
		e := Images("li img", "src")

		val, err := e.Extract(Fragment(`<ul><li><img src="/a.png"></li><li><img src="b.png " width="500" height="900"></li></ul>`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal([]*ScrapedArticleImage{
			{Url: "/a.png"},
			{Url: "b.png", Width: 500, Height: 900},
		}))

		val, err = e.Extract(Fragment(`<img src="foo.jpg">`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeNil())

		e = Images("li a", "href")

		val, err = e.Extract(Fragment(`<ul><li><a href="/a.png"></li><li><a href="b.png"></li></ul>`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal([]*ScrapedArticleImage{
			{Url: "/a.png"},
			{Url: "b.png"},
		}))
	})
})
