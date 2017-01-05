package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega/gexec"
)

var (
	pathToNsProcessCLI string
	rootfsFilepath     string
)

func TestNsProcess(t *testing.T) {
	BeforeSuite(func() {
		var err error

		// build the ns-process binary
		pathToNsProcessCLI, err = gexec.Build("github.com/teddyking/ns-process")
		Expect(err).NotTo(HaveOccurred())

		// setup a test rootfs dir
		busyboxTarFilepath := filepath.Join("assets", "busybox.tar")
		rootfsFilepath = filepath.Join("/tmp", "ns-process", "test", "rootfs")

		Expect(os.RemoveAll(rootfsFilepath)).To(Succeed())
		Expect(os.MkdirAll(rootfsFilepath, 0755)).To(Succeed())

		untar(busyboxTarFilepath, rootfsFilepath)

		// check for netsetgo
		checkNetsetgo()
	})

	AfterEach(func() {
		// allow time for various network operations to complete between tests
		// super gross hack I know, but the tests here are very simplified and
		// really serve only to support a basic TDD workflow
		time.Sleep(time.Second / 2)
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

	SetDefaultEventuallyTimeout(time.Second * 5)
	RegisterFailHandler(Fail)
	RunSpecs(t, "ns-process")
}

func inode(pid, namespaceType string) string {
	namespace, err := os.Readlink(fmt.Sprintf("/proc/%s/ns/%s", pid, namespaceType))
	Expect(err).NotTo(HaveOccurred())

	requiredFormat := regexp.MustCompile(`^\w+:\[\d+\]$`)
	Expect(requiredFormat.MatchString(namespace)).To(BeTrue())

	namespace = strings.Split(namespace, ":")[1]
	namespace = namespace[1:]
	namespace = namespace[:len(namespace)-1]

	return namespace
}

func untar(src, dst string) {
	// use tar command for sake of simplicity
	untarCmd := exec.Command("tar", "-C", dst, "-xf", src)
	Expect(untarCmd.Run()).To(Succeed())
}

func checkNetsetgo() {
	_, err := exec.LookPath("netsetgo")
	Expect(err).NotTo(HaveOccurred())
}
