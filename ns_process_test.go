package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os/exec"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("The ns-process CLI", func() {
	var (
		session *gexec.Session
	)

	BeforeEach(func() {
		var err error

		command := exec.Command(pathToNsProcessCLI)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	It("exits with a 0 exit code", func() {
		Eventually(session).Should(gexec.Exit(0))
	})
})
