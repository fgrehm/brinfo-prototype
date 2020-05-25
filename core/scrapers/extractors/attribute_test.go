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

	It("can be marked as optional", func() {
		e := OptAttribute("meta", "name")

		val, err := e.Extract(Fragment(`<meta name="bla">`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(Equal("bla"))

		val, err = e.Extract(Fragment(`<metas name="bla">`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeNil())

		val, err = e.Extract(Fragment(`<meta names="bla">`))
		Expect(err).NotTo(HaveOccurred())
		Expect(val).To(BeNil())

		_, err = e.Extract(Fragment(`<meta name="bla"><meta name="bla2">`))
		Expect(err).To(HaveOccurred())
	})
})
