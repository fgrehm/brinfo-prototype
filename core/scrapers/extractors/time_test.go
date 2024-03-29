package extractors_test

import (
	"time"

	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Date extractors", func() {
	brLoc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	Describe("TimeText", func() {
		It("parses dates formatted in day/month", func() {
			e := TimeText("span")

			val, err := extract(e, `<p>Publicado em <span>22/02/2020 15:50</span></p>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2020, 2, 22, 15, 50, 0, 0, brLoc)))

			val, err = extract(e, `<p>Publicado em <span>21/02/2020 16h50</span></p>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2020, 2, 21, 16, 50, 0, 0, brLoc)))

			val, err = extract(e, `<span>21/02/2020 - 16:50</span>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2020, 2, 21, 16, 50, 0, 0, brLoc)))

			val, err = extract(e, `<span>21/02/2020 - 16h50</span>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2020, 2, 21, 16, 50, 0, 0, brLoc)))

			val, err = extract(e, `<span>21/02/2020 - 16:50:11</span>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2020, 2, 21, 16, 50, 11, 0, brLoc)))

			e = TimeText("p")
			val, err = extract(e, `<p>21 de FeVereiro dE 2020 àS 16:26 | <span>Geral</span></p>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2020, 2, 21, 16, 26, 0, 0, brLoc)))

			e = TimeText("p")
			val, err = extract(e, `<p>21/fevereiro/2020 4:26 pm | <span>Geral</span></p>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2020, 2, 21, 16, 26, 0, 0, brLoc)))

			val, err = extract(e, `<p><strong>Publicação</strong>: 22 de março dE 2020 àS 16:22 | <span>Geral</span></p>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2020, 3, 22, 16, 22, 0, 0, brLoc)))
		})

		It("errors if element not found", func() {
			e := TimeText("span")

			val, err := extract(e, `<p>Publicado em <time>22/02/2020 15:50</time></p>`)
			Expect(err).To(HaveOccurred())
			Expect(val).To(BeNil())
		})

		It("can be made optional", func() {
			e := OptTimeText("span")

			val, err := extract(e, `<p>Publicado em <time>22/02/2020 15:50</time></p>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())
		})
	})

	Describe("TimeAttribute", func() {
		It("parses json like timestamps", func() {
			e := TimeAttribute("time", "datetime")

			val, err := extract(e, `<p><time datetime="2010-02-21 15:50:00">foobar</time></p>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2010, 2, 21, 15, 50, 0, 0, brLoc)))
		})

		It("errors if attribute not found", func() {
			e := TimeAttribute("time", "datetime")

			val, err := extract(e, `<p><time pubdates="2010-02-21 15:50:00">foobar</time></p>`)
			Expect(err).To(HaveOccurred())
			Expect(val).To(BeNil())
		})

		It("can be made optional", func() {
			e := OptTimeAttribute("time", "datetime")

			val, err := extract(e, `<p><time pubdate="2010-02-21 15:50:00">foobar</time></p>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())

			val, err = extract(e, `<p><time datetime="2010-02-21 15:50:00">foobar</time></p>`)
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())
			Expect(val).To(Equal(time.Date(2010, 2, 21, 15, 50, 0, 0, brLoc)))
		})
	})
})
