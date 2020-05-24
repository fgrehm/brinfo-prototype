package extractors_test

import (
	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Text", func() {
	It("works as expected", func() {
		val, err := Text("h1.title", false).Extract(Fragment(`<h1 class="title">foo</h1>`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("foo"))

		val, err = Text("h1.title", true).Extract(Fragment(`<h1 class="title">foo</h1>`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal([]string{"foo"}))

		_, err = Text("p", false).Extract(Fragment(`<p>a</p><i>b2</i><p>b</p><i>a2</i>`))
		Expect(err).To(HaveOccurred())

		val, err = Text("p", true).Extract(Fragment(`<p>a</p><i>b2</i><p>b</p><i>a2</i>`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal([]string{"a", "b"}))

		_, err = Text("h1.title", false).Extract(Fragment(`<h1 class="titles">foo</h1>`))
		Expect(err).To(HaveOccurred())

		_, err = Text("h1.title", true).Extract(Fragment(`<h1 class="titles">foo</h1>`))
		Expect(err).To(HaveOccurred())
	})

	// 	XIt("can be made optional with ExtractOptionalText", func() {
})
