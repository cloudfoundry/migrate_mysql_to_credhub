package main_test

import (
	"errors"
	"os/exec"

	. "code.cloudfoundry.org/migrate_mysql_to_credhub"

	"github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	Describe("#Main", func() {
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

	Describe("#HandleSQLStoreError", func() {
		Context("when given a nil error", func() {
			It("should return nil", func() {
				Expect(HandleSQLStoreError(nil)).To(BeNil())
			})
		})

		Context("when given a generic error", func() {
			It("should return the error", func() {
				err := errors.New("uhhoh")
				Expect(HandleSQLStoreError(err)).To(MatchError(err))
			})
		})

		Context("when given a MySQL error #1049", func() {
			It("should return nil", func() {
				err := &mysql.MySQLError{
					Number: 1049,
				}
				Expect(HandleSQLStoreError(err)).To(BeNil())
			})
		})
	})
})
