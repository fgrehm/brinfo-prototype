package testutils_test

import (
	"testing"

	_ "github.com/fgrehm/brinfo/core/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "testutils Suite")
}
