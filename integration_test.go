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

		It("starts /bin/sh in a new set of namespaces", func() {
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
				cmdToRunInNamespacedShell = "ls && pwd && mount"
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

			It("mounts /proc", func() {
				Eventually(stdout).Should(gbytes.Say("proc on /proc type proc"))
			})
		})

		Describe("network namespace", func() {
			BeforeEach(func() {
				// 10.10.10.1 is the default IP address of the bridge assigned by netsetgo
				cmdToRunInNamespacedShell = "ip addr && ping -c 1 -W 1 10.10.10.1"
			})

			It("assigns an interface/IP address", func() {
				// 10.10.10.2 is the default IP address assigned by netsetgo
				Eventually(stdout).Should(gbytes.Say("inet 10.10.10.2/24 scope global veth1"))
			})

			It("is able to ping the gateway", func() {
				Eventually(stdout).Should(gbytes.Say("1 packets transmitted, 1 packets received, 0% packet loss"))
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
An example rootfs, BusyBox, can be downloaded and unpacked as follows:

wget "https://raw.githubusercontent.com/teddyking/ns-process/4.0/assets/busybox.tar"
mkdir -p /does/not/exist
tar -C /does/not/exist -xf busybox.tar
`
			Eventually(stdout).Should(gbytes.Say(usefulErrorMsg))
		})
	})

	Context("when netsetgo does not exist", func() {
		BeforeEach(func() {
			args = append(args, "-netsetgo", "/does/not/exist")
		})

		It("provides a helpful error message", func() {
			usefulErrorMsg := `
Unable to find the netsetgo binary at "/does/not/exist".
netsetgo is an external binary used to configure networking.
You must download netsetgo, chown it to the root user and apply the setuid bit.
This can be done as follows:

wget "https://github.com/teddyking/netsetgo/releases/download/0.0.1/netsetgo"
sudo mv netsetgo /usr/local/bin/
sudo chown root:root /usr/local/bin/netsetgo
sudo chmod 4755 /usr/local/bin/netsetgo
`
			Eventually(stdout).Should(gbytes.Say(usefulErrorMsg))
		})
	})
})
