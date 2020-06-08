package extractors_test

import (
	"time"

	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FromString", func() {
	brLoc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	Context("single attribute extraction", func() {
		It("works for required data", func () {
			e, err := FromString("p a | href")
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())

			val, err := e.Extract(Fragment(`<p>other</p><p>Foo <a href="#foo">link</a></p><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal("#foo"))

			val, err = e.Extract(Fragment(`<p>other</p><div>Foo <a href="#foo">link</a></div><a href="#bla">foo</a>`))
			Expect(err).To(HaveOccurred())
		})

		It("works for optional data", func () {
			e, err := FromString("p a | href?")
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())

			val, err := e.Extract(Fragment(`<p>other</p><p>Foo <a href="#foo">link</a></p><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal("#foo"))

			val, err = e.Extract(Fragment(`<p>other</p><div>Foo <a href="#foo">link</a></div><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())
		})
	})

	Context("single text extraction", func() {
		It("works for required data", func () {
			e, err := FromString("p a | text")
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())

			val, err := e.Extract(Fragment(`<p>other</p><p>Foo <a href="#foo">link</a></p><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal("link"))

			val, err = e.Extract(Fragment(`<p>other</p><div>Foo <a href="#foo">link</a></div><a href="#bla">foo</a>`))
			Expect(err).To(HaveOccurred())
		})

		It("works for optional data", func () {
			e, err := FromString("p a | text?")
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())

			val, err := e.Extract(Fragment(`<p>other</p><p>Foo <a href="#foo">link</a></p><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal("link"))

			val, err = e.Extract(Fragment(`<p>other</p><div>Foo <a href="#foo">link</a></div><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())
		})
	})

	Context("time attribute extraction", func() {
		It("works for required data", func () {
			e, err := FromString("time.pub | pubdate::time")
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())

			val, err := e.Extract(Fragment(`<p>other</p><p>Foo <time class="pub" pubdate="2020-03-20 18:30:00">whatever</time></p><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(time.Date(2020, 3, 20, 18, 30, 0, 0, brLoc)))

			val, err = e.Extract(Fragment(`<p>other</p><p>Foo <time pubdate="2020-03-20 18:30:00">whatever</time></p><a href="#bla">foo</a>`))
			Expect(err).To(HaveOccurred())
		})

		It("works for optional data", func () {
			e, err := FromString("time.pub | pubdate?::time")
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())

			val, err := e.Extract(Fragment(`<p>other</p><p>Foo <time class="pub" pubdate="2020-03-20 18:30:00">whatever</time></p><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(time.Date(2020, 3, 20, 18, 30, 0, 0, brLoc)))

			val, err = e.Extract(Fragment(`<p>other</p><p>Foo <time pubdate="2020-03-20 18:30:00">whatever</time></p><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("time text extraction", func() {
		It("works for required data", func () {
			e, err := FromString("p em | text::time")
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())

			val, err := e.Extract(Fragment(`<p>other</p><p>Foo <em>20/03/2020 18:30</em></p><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2020, 3, 20, 18, 30, 0, 0, brLoc)))

			val, err = e.Extract(Fragment(`<p>other</p><div>Foo <em>20/03/2020 18:30</em></div><a href="#bla">foo</a>`))
			Expect(err).To(HaveOccurred())
		})

		It("works for optional data", func () {
			e, err := FromString("p em | text?::time")
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())

			val, err := e.Extract(Fragment(`<p>other</p><p>Foo <em>20/03/2020 18:30</em></p><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(time.Date(2020, 3, 20, 18, 30, 0, 0, brLoc)))

			val, err = e.Extract(Fragment(`<p>other</p><div>Foo <em>20/03/2020 18:30</em></div><a href="#bla">foo</a>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())
		})
	})
})
