package extractors_test

import (
	. "github.com/fgrehm/brinfo/core/scrapers/extractors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("BasicArticle", func() {
	It("works when og tags are present", func() {
		e := BasicArticle()

		val, err := extract(e, basicArticleWithOGHTML)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).NotTo(BeNil())
		Expect(val).To(MatchKeys(IgnoreExtras, Keys{
			"title":       Equal("Article title"),
			"excerpt":     Equal("From og meta with a few words"),
			"fullText":    Equal("Article body"),
			"publishedAt": Not(BeNil()),
			"modifiedAt":  Not(BeNil()),
			"imageURL":    Equal("https://image.url"),
		}))
	})

	It("works when meta tags are present", func() {
		e := BasicArticle()

		val, err := extract(e, basicArticleWithMetaHTML)
		Expect(err).NotTo(HaveOccurred())
		Expect(val).NotTo(BeNil())
		Expect(val).To(MatchKeys(IgnoreExtras, Keys{
			"title":       Equal("Article title | Whatever"),
			"excerpt":     Equal("From meta with a few words"),
			"fullText":    Equal("Other article body"),
			"publishedAt": Not(BeNil()),
			"modifiedAt":  Not(BeNil()),
			"imageURL":    Equal("https://image.url"),
		}))
	})
})

var basicArticleWithOGHTML = `<html>
	<head>
		<title>Article title - Website</title>
		<meta property="og:site_name" content="Website">
		<meta property="og:description" content="From og meta with a few words">
		<meta property="og:type" content="article">
		<meta property="og:image" content="https://image.url">
		<meta property="article:published_time" content="2020-06-21T15:53:10-03:00">
		<meta property="article:modified_time" content="2020-06-21T16:52:10-03:00">
	</head>
	<body>
		<p>Article body</p>
	</body>
</html>`

var basicArticleWithMetaHTML = `<html>
	<head>
		<title>Article title | Whatever | Website</title>
		<meta name="description" content="From meta with a few words">
		<meta property="article:published_time" content="2020-06-21 15:53:10">
		<meta property="article:modified_time" content="2020-06-21 16:52:10">
	</head>
	<body>
		<p>
			<img src="https://image.url">
			Other article body
		</p>
	</body>
</html>`
