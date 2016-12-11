package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"regexp"
	"strings"
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
