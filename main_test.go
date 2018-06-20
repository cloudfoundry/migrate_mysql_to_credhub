package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	It("runs with specified arguments", func() {
		args := []string{
			"--dbDriver", "mysql",
			"--dbUsername", "some-db-username",
			"--dbPassword", "some-db-password",
			"--dbHostname", "some-db-hostname",
			"--dbPort", "1234",
			"--dbName", "some-db-name",
			"--credhubURL", "some-credhub-url",
			"--storeID", "some-store-id",
			"--uaaClientID", "some-uaa-client-id",
			"--uaaClientSecret", "some-uaa-client-secret",
		}
		session, err := gexec.Start(exec.Command(binaryPath, args...), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		<-session.Exited
		Expect(session.ExitCode()).To(Equal(0))
	})

	It("fails if required argument is not provided", func() {
		args := []string{
			"--dbUsername", "some-db-username",
			"--dbPassword", "some-db-password",
			"--dbHostname", "some-db-hostname",
			"--dbPort", "1234",
			"--dbName", "some-db-name",
			"--credhubURL", "some-credhub-url",
			"--storeID", "some-store-id",
			"--uaaClientID", "some-uaa-client-id",
			"--uaaClientSecret", "some-uaa-client-secret",
		}
		session, err := gexec.Start(exec.Command(binaryPath, args...), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		<-session.Exited
		Expect(session.ExitCode()).NotTo(Equal(0))
		Expect(session.Err).Should(Say("dbDriver"))
	})
})
