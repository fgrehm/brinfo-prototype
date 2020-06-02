package inmemory_test

import (
	"context"

	. "github.com/fgrehm/brinfo/core"
	mem "github.com/fgrehm/brinfo/storage/inmemory"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContentSource", func() {
	It("maintains a registry of content sources", func() {
		ctx := context.Background()
		r := mem.NewContentSourceRepo()

		cs := &ContentSource{ID: "br-foo", Host: "example.com"}
		err := r.Register(cs)
		Expect(err).NotTo(HaveOccurred())

		// No duplicates
		err = r.Register(&ContentSource{ID: "br-foo"})
		Expect(err).To(HaveOccurred())

		// Lookup by guid
		cs, err = r.FindByID(ctx, cs.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(cs).To(Equal(cs))

		_, err = r.FindByID(ctx, "NO")
		Expect(err).To(HaveOccurred())

		// Lookup by host
		cs, err = r.FindByHost(ctx, cs.Host)
		Expect(err).NotTo(HaveOccurred())
		Expect(cs).To(Equal(cs))

		_, err = r.FindByHost(ctx, "NO")
		Expect(err).To(HaveOccurred())
	})
})
