package extractors_test

import (
	"time"

	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Loaders", func() {
	brLoc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}
	Describe("FromString", func() {
		Context("single attribute extraction", func() {
			It("works for required data", func() {
				e, err := FromString("p a | href")
				Expect(err).NotTo(HaveOccurred())
				Expect(e).NotTo(BeNil())

				val, err := extract(e, `<p>other</p><p>Foo <a href="#foo">link</a></p><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("#foo"))

				_, err = extract(e, `<p>other</p><div>Foo <a href="#foo">link</a></div><a href="#bla">foo</a>`)
				Expect(err).To(HaveOccurred())
			})

			It("works for optional data", func() {
				e, err := FromString("p a | href?")
				Expect(err).NotTo(HaveOccurred())
				Expect(e).NotTo(BeNil())

				val, err := extract(e, `<p>other</p><p>Foo <a href="#foo">link</a></p><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("#foo"))

				val, err = extract(e, `<p>other</p><div>Foo <a href="#foo">link</a></div><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(BeNil())
			})
		})

		Context("single text extraction", func() {
			It("works for required data", func() {
				e, err := FromString("p a | text")
				Expect(err).NotTo(HaveOccurred())
				Expect(e).NotTo(BeNil())

				val, err := extract(e, `<p>other</p><p>Foo <a href="#foo">link</a></p><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("link"))

				_, err = extract(e, `<p>other</p><div>Foo <a href="#foo">link</a></div><a href="#bla">foo</a>`)
				Expect(err).To(HaveOccurred())
			})

			It("works for optional data", func() {
				e, err := FromString("p a | text?")
				Expect(err).NotTo(HaveOccurred())
				Expect(e).NotTo(BeNil())

				val, err := extract(e, `<p>other</p><p>Foo <a href="#foo">link</a></p><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("link"))

				val, err = extract(e, `<p>other</p><div>Foo <a href="#foo">link</a></div><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(BeNil())
			})
		})

		Context("time attribute extraction", func() {
			It("works for required data", func() {
				e, err := FromString("time.pub | pubdate::time")
				Expect(err).NotTo(HaveOccurred())
				Expect(e).NotTo(BeNil())

				val, err := extract(e, `<p>other</p><p>Foo <time class="pub" pubdate="2020-03-20 18:30:00">whatever</time></p><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal(time.Date(2020, 3, 20, 18, 30, 0, 0, brLoc)))

				_, err = extract(e, `<p>other</p><p>Foo <time pubdate="2020-03-20 18:30:00">whatever</time></p><a href="#bla">foo</a>`)
				Expect(err).To(HaveOccurred())
			})

			It("works for optional data", func() {
				e, err := FromString("time.pub | pubdate?::time")
				Expect(err).NotTo(HaveOccurred())
				Expect(e).NotTo(BeNil())

				val, err := extract(e, `<p>other</p><p>Foo <time class="pub" pubdate="2020-03-20 18:30:00">whatever</time></p><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal(time.Date(2020, 3, 20, 18, 30, 0, 0, brLoc)))

				_, err = extract(e, `<p>other</p><p>Foo <time pubdate="2020-03-20 18:30:00">whatever</time></p><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("time text extraction", func() {
			It("works for required data", func() {
				e, err := FromString("p em | text::time")
				Expect(err).NotTo(HaveOccurred())
				Expect(e).NotTo(BeNil())

				val, err := extract(e, `<p>other</p><p>Foo <em>20/03/2020 18:30</em></p><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).NotTo(BeNil())
				Expect(val).To(Equal(time.Date(2020, 3, 20, 18, 30, 0, 0, brLoc)))

				_, err = extract(e, `<p>other</p><div>Foo <em>20/03/2020 18:30</em></div><a href="#bla">foo</a>`)
				Expect(err).To(HaveOccurred())
			})

			It("works for optional data", func() {
				e, err := FromString("p em | text?::time")
				Expect(err).NotTo(HaveOccurred())
				Expect(e).NotTo(BeNil())

				val, err := extract(e, `<p>other</p><p>Foo <em>20/03/2020 18:30</em></p><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal(time.Date(2020, 3, 20, 18, 30, 0, 0, brLoc)))

				val, err = extract(e, `<p>other</p><div>Foo <em>20/03/2020 18:30</em></div><a href="#bla">foo</a>`)
				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(BeNil())
			})
		})
	})

	Describe("FromJSON", func() {
		It("works for top level fields", func() {
			e, err := FromJSON([]byte(`{"title": "p a | href"}`))
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(len(e)).To(Equal(1))

			val, err := extract(e[0], `<html><body><p>other</p><p>Foo <a href="#foo">link</a></p><a href="#bla">foo</a></body></html>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(map[string]ExtractorResult{"title": "#foo"}))

			_, err = extract(e[0], `<html><body><p>other</p><div>Foo <a href="#foo">link</a></div><a href="#bla">foo</a></body></html>`)
			Expect(err).To(HaveOccurred())
		})

		// Temporarily disabled while we make it deterministic, we can't rely on map keys ordering
		XIt("works when wrappers are specified", func() {
			e, err := FromJSON([]byte(`{"head": { "title": "title | text" }, "body": {"fullText": "#main | text"}}`))
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(len(e)).To(Equal(2))

			values := []ExtractorResult{}

			val, err := extract(e[0], `<html><head><title>Page title</title></head><body><title>bla</title></body></html>`)
			Expect(err).NotTo(HaveOccurred())
			values = append(values, val)

			val, err = extract(e[1], `<html><head><title id="main">Page title</title></head><body><div id="main">Contents</div></body></html>`)
			Expect(err).NotTo(HaveOccurred())
			values = append(values, val)

			Expect(values).To(ContainElement(map[string]ExtractorResult{"title": "Page title"}))
			Expect(values).To(ContainElement(map[string]ExtractorResult{"fullText": "Contents"}))
		})

		It("normalizes attributes", func() {
			e, err := FromJSON([]byte(`{"full_text": "p | text"}`))
			Expect(err).NotTo(HaveOccurred())
			Expect(e).NotTo(BeNil())
			Expect(len(e)).To(Equal(1))

			val, err := extract(e[0], `<html><body><p>full text</p></body></html>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(map[string]ExtractorResult{"fullText": "full text"}))
		})

		It("errors if can't parse extractors", func() {
			e, err := FromJSON([]byte(`{"full_text": "p"}`))
			Expect(err).To(HaveOccurred())
			Expect(e).To(BeNil())

			e, err = FromJSON([]byte(`{"head": {"full_text": "p"}}`))
			Expect(err).To(HaveOccurred())
			Expect(e).To(BeNil())
		})
	})
})
