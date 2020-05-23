package cmd

import (
	"testing"

	op "github.com/fgrehm/brinfo/core/operations"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCore(t *testing.T) {
	op.UseCache = false

	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI commands")
	// Most of the logic should be self contained in the core pkg, this suite
	// exists only to compile the CLI as part of a test run
}
