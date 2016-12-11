package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bufio"
	"io"
	"os/exec"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("The ns-process CLI", func() {
	var (
		command *exec.Cmd
		session *gexec.Session
		stdin   io.WriteCloser
		stdout  *gbytes.Buffer
	)

	BeforeEach(func() {
		var err error

		command = exec.Command(pathToNsProcessCLI)
		stdin, err = command.StdinPipe()
		Expect(err).NotTo(HaveOccurred())
		stdout = gbytes.NewBuffer()

		stdinWriter := bufio.NewWriter(stdin)
		stdinWriter.WriteString("readlink /proc/self/ns/uts")
		stdinWriter.Flush()
		Expect(stdin.Close()).To(Succeed())
	})

	JustBeforeEach(func() {
		var err error

		session, err = gexec.Start(command, stdout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(stdout.Close()).To(Succeed())
	})

	It("exits with a 0 exit code", func() {
		Eventually(session).Should(gexec.Exit(0))
	})

	It("starts a /bin/sh process in a new UTS namespace", func() {
		currentUTSInode := inode("self", "uts")
		Consistently(stdout).ShouldNot(gbytes.Say(currentUTSInode))
	})
})
