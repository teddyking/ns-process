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
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
	})

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
