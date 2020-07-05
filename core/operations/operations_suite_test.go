package operations_test

import (
	"context"
	"testing"

	. "github.com/fgrehm/brinfo/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOperations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Operations Suite")
}

type fakeScraper struct {
	data *ArticleData
}

func (f *fakeScraper) Run(context.Context, []byte, string, string) (*ArticleData, error) {
	return f.data, nil
}
