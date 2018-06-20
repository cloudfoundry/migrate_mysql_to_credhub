package main_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var binaryPath string

func TestMigrateMysqlToCredhub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MigrateMysqlToCredhub Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	var err error
	binaryPath, err = gexec.Build("code.cloudfoundry.org/migrate_mysql_to_credhub", "-race")
	Expect(err).NotTo(HaveOccurred())

	return []byte(binaryPath)
}, func(bytes []byte) {
	binaryPath = string(bytes)
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	err := os.RemoveAll(binaryPath)
	Expect(err).NotTo(HaveOccurred())
})
