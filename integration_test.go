package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("ns-process", func() {
	var (
		command                   *exec.Cmd
		args                      []string
		session                   *gexec.Session
		stdin                     io.WriteCloser
		cmdToRunInNamespacedShell string
		stdout                    *gbytes.Buffer
	)

	BeforeEach(func() {
		args = []string{"-rootfs", rootfsFilepath}
		cmdToRunInNamespacedShell = "true"
	})

	JustBeforeEach(func() {
		var err error

		command = exec.Command(pathToNsProcessCLI, args...)
		stdin, err = command.StdinPipe()
		Expect(err).NotTo(HaveOccurred())
		stdout = gbytes.NewBuffer()

		stdinWriter := bufio.NewWriter(stdin)
		stdinWriter.WriteString(cmdToRunInNamespacedShell)
		stdinWriter.Flush()
		Expect(stdin.Close()).To(Succeed())

		session, err = gexec.Start(command, stdout, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(stdout.Close()).To(Succeed())
	})

	Describe("cloning namespaces", func() {
		BeforeEach(func() {
			cmdToRunInNamespacedShell = "ls -lah /proc/self/ns"
		})

		// TODO: unpend this test once we have a /proc again
		PIt("starts /bin/sh in a new set of namespaces", func() {
			namespaces := []string{"mnt", "uts", "ipc", "pid", "net", "user"}

			for _, namespace := range namespaces {
				currentInode := inode("self", namespace)
				Consistently(stdout).ShouldNot(gbytes.Say(currentInode))
			}
		})
	})

	Describe("namespace setup and configuration", func() {
		Describe("user namespace", func() {
			BeforeEach(func() {
				cmdToRunInNamespacedShell = "id"
			})

			It("applies a UID mapping", func() {
				Eventually(stdout).Should(gbytes.Say(`uid=0\(root\)`))
			})

			It("applies a GID mapping", func() {
				Eventually(stdout).Should(gbytes.Say(`gid=0\(root\)`))
			})
		})

		Describe("mount namespace", func() {
			BeforeEach(func() {
				_, err := os.Create(filepath.Join(rootfsFilepath, "in-newroot"))
				Expect(err).NotTo(HaveOccurred())
				cmdToRunInNamespacedShell = "ls && pwd"
			})

			AfterEach(func() {
				Expect(os.Remove(filepath.Join(rootfsFilepath, "in-newroot"))).To(Succeed())
			})

			It("starts /bin/sh with a new root filesystem", func() {
				Eventually(stdout).Should(gbytes.Say("in-newroot"))
			})

			It("sets / for /bin/sh to the new root filesystem", func() {
				Consistently(stdout).ShouldNot(gbytes.Say("/.pivot_root"))
				Eventually(stdout).Should(gbytes.Say("/"))
			})
		})
	})

	Context("when the rootfs directory does not exist", func() {
		BeforeEach(func() {
			args = []string{"-rootfs", "/does/not/exist"}
		})

		It("exits with a 1 exit code", func() {
			Eventually(session).Should(gexec.Exit(1))
		})

		It("provides a helpful error message", func() {
			usefulErrorMsg := `
"/does/not/exist" does not exist.
Please create this directory and unpack a suitable root filesystem inside it.
An example rootfs, BusyBox, can be downloaded from:

https://raw.githubusercontent.com/teddyking/ns-process/4.0/assets/busybox.tar

And unpacked by:

mkdir -p /does/not/exist
tar -C /does/not/exist -xf busybox.tar
`
			Eventually(stdout).Should(gbytes.Say(usefulErrorMsg))
		})
	})
})
