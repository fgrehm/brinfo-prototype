package extractors_test

import (
	"time"

	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PublishedDates", func() {
	brLoc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	Context("meta[article:*]", func() {
		It("extracts published_at from article:published_time", func() {
			e := PublishedDates()
			val, err := e.Extract(Fragment(`<head><meta property="article:published_time" content="2010-02-21 15:50"></head> `))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())

			data, ok := val.(map[string]*time.Time)
			if !ok {
				panic("Returned something weird")
			}
			Expect(data["published_at"]).NotTo(BeNil())
			Expect(*data["published_at"]).To(Equal(time.Date(2010, 2, 21, 15, 50, 0, 0, brLoc)))
			Expect(data["modified_at"]).To(BeNil())
		})

		It("extracts updated_at from article:modified_time", func() {
			e := PublishedDates()
			val, err := e.Extract(Fragment(`<head><meta property="article:modified_time" content="2010-02-21 15:50"></head> `))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())

			data, ok := val.(map[string]*time.Time)
			if !ok {
				panic("Returned something weird")
			}
			Expect(data["modified_at"]).NotTo(BeNil())
			Expect(*data["modified_at"]).To(Equal(time.Date(2010, 2, 21, 15, 50, 0, 0, brLoc)))
			Expect(data["published_at"]).To(BeNil())
		})

		It("is restricted to elements within <head>", func() {
			e := PublishedDates()
			val, err := e.Extract(Fragment(`<body><meta property="article:published_time" content="2010-02-21 15:50"></body> `))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())
		})
	})

	Context("article time", func() {
		It("extracts published_at from <article><time> with pubdate", func() {
			e := PublishedDates()

			val, err := e.Extract(Fragment(`<article>
				<time pubdate="" datetime="2010-02-21 15:50:00 -0300">foobar</time></article>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())

			data, ok := val.(map[string]*time.Time)
			if !ok {
				panic("Returned something weird")
			}
			Expect(data["published_at"]).NotTo(BeNil())
			Expect(*data["published_at"]).To(Equal(time.Date(2010, 2, 21, 15, 50, 0, 0, brLoc)))
			Expect(data["modified_at"]).To(BeNil())
		})

		It("is restricted to elements with proper attrs", func() {
			e := PublishedDates()

			val, err := e.Extract(Fragment(`<article>
				<time pubdates="" datetime="2010-02-21 15:50:00 -0300">foobar</time></article>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())

			val, err = e.Extract(Fragment(`<article>
				<time pubdate="" datetimes="2010-02-21 15:50:00 -0300">foobar</time></article>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())
		})
	})

	Context("rnews", func() {
		It("extracts published_at from rnews:datePublished", func() {
			e := PublishedDates()

			val, err := e.Extract(Fragment(`<body>
				<article vocab="http://schema.org/" typeof="Article" prefix="rnews: http://iptc.org/std/rNews/2011-10-07#">
					<span class="documentPublished">
						<span>publicado</span>:
						<span property="rnews:datePublished">21/02/2010 15h50</span>,
					</span>
				</article>
			</body>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())

			data, ok := val.(map[string]*time.Time)
			if !ok {
				panic("Returned something weird")
			}
			Expect(data["published_at"]).NotTo(BeNil())
			Expect(*data["published_at"]).To(Equal(time.Date(2010, 2, 21, 15, 50, 0, 0, brLoc)))
			Expect(data["modified_at"]).To(BeNil())
		})

		It("extracts updated_at from rnews:dateModified", func() {
			e := PublishedDates()

			val, err := e.Extract(Fragment(`<body>
				<div vocab="http://schema.org/" typeof="Article" prefix="rnews: http://iptc.org/std/rNews/2011-10-07#">
					<span class="documentPublished">
						<span>atualizado</span>:
						<span property="rnews:dateModified"> 21/02/2010 15h50 </span>,
					</span>
				</div>
			</body>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).NotTo(BeNil())

			data, ok := val.(map[string]*time.Time)
			if !ok {
				panic("Returned something weird")
			}
			Expect(data["modified_at"]).NotTo(BeNil())
			Expect(*data["modified_at"]).To(Equal(time.Date(2010, 2, 21, 15, 50, 0, 0, brLoc)))
			Expect(data["published_at"]).To(BeNil())
		})

		It("is restricted to elements with proper attrs", func() {
			e := PublishedDates()

			val, err := e.Extract(Fragment(`<body>
				<div vocab="http://schemas.org/" typeof="Article" prefix="rnews: http://iptc.org/std/rNews/2011-10-07#">
					<span property="rsnews:dateModified">21/02/2010 15h50</span>
				</div></body>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())

			val, err = e.Extract(Fragment(`<body>
				<div vocab="http://schema.org/" typeof="Articles" prefix="rnews: http://iptc.org/std/rNews/2011-10-07#">
					<span property="rnews:dasteModified">21/02/2010 15h50</span>
				</div></body>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())

			val, err = e.Extract(Fragment(`<body>
				<div vocab="http://schema.org/" typeof="Article" prefix="rsnews: http://iptc.org/std/rNews/2011-10-07#">
					<span property="rnews:dateModified">21/02/2010 15h50</span>
				</div></body>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())

			val, err = e.Extract(Fragment(`<body>
				<div vocab="http://schema.org/" typeof="Article" prefix="rnews: http://iptc.org/std/rNews/2011-10-07#">
					<span other="rnews:dateModified">21/02/2010 15h50</span>
				</div></body>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())

			val, err = e.Extract(Fragment(`<body>
				<div vocab="http://schema.org/" typeof="Article" prefix="rnews: http://iptc.org/std/rNews/2011-10-07#">
					<span property="rnesws:dateModified">21/02/2010 15h50</span>
				</div></body>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())

			val, err = e.Extract(Fragment(`<body>
				<div vocab="http://schema.org/" typeof="Article" prefix="rnews: http://iptc.org/std/rNews/2011-10-07#">
					<span property="rnews:datePodified">21/02/2010 15h50</span>
				</div></body>`))
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(BeNil())
		})
	})
})
