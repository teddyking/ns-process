package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/onsi/gomega/gexec"
)

var pathToNsProcessCLI string

func TestNsProcess(t *testing.T) {
	BeforeSuite(func() {
		var err error

		pathToNsProcessCLI, err = gexec.Build("github.com/teddyking/ns-process")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "ns-process")
}
