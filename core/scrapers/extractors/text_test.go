package extractors_test

import (
	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Text", func() {
	It("works as expected", func() {
		val, err := extract(Text("h1.title", false), `<h1 class="title">  foo  </h1>`)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("foo"))

		val, err = extract(Text("h1.title", true), `<h1 class="title">foo</h1>`)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal([]string{"foo"}))

		_, err = extract(Text("p", false), `<p>a</p><i>b2</i><p>b</p><i>a2</i>`)
		Expect(err).To(HaveOccurred())

		val, err = extract(Text("p", true), `<p>a</p><i>b2</i><p>b</p><i>a2</i>`)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal([]string{"a", "b"}))

		_, err = extract(Text("h1.title", false), `<h1 class="titles">foo</h1>`)
		Expect(err).To(HaveOccurred())

		_, err = extract(Text("h1.title", true), `<h1 class="titles">foo</h1>`)
		Expect(err).To(HaveOccurred())
	})

	It("can be marked as optional", func() {
		e := OptText("p", false)

		val, err := extract(e, `<p>bla</p>`)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("bla"))

		val, err = extract(e, `<div>a<div>`)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeNil())

		val, err = extract(e, `<p></p>`)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeNil())

		_, err = extract(e, `<p>a</p><p>b</p>`)
		Expect(err).To(HaveOccurred())
	})
})
