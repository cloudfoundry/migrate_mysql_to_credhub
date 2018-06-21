package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/brokerapi"

	"code.cloudfoundry.org/lager/lagertest"
	. "code.cloudfoundry.org/migrate_mysql_to_credhub"
	"code.cloudfoundry.org/service-broker-store/brokerstore"
	"code.cloudfoundry.org/service-broker-store/brokerstore/brokerstorefakes"
)

var _ = Describe("Migrator", func() {
	var (
		migrator  Migrator
		fromStore *brokerstorefakes.FakeStore
		toStore   *brokerstorefakes.FakeStore
	)

	BeforeEach(func() {
		logger := lagertest.NewTestLogger("migrator-test")
		migrator = NewMigrator(logger)
		fromStore = &brokerstorefakes.FakeStore{}
		toStore = &brokerstorefakes.FakeStore{}
	})

	Context("when there are instance details in fromStore", func() {
		BeforeEach(func() {
			fromStore.RetrieveAllInstanceDetailsReturns(map[string]brokerstore.ServiceInstance{
				"123": brokerstore.ServiceInstance{ServiceID: "some-service-1"},
				"456": brokerstore.ServiceInstance{ServiceID: "some-service-2"},
			}, nil)
			err := migrator.Migrate(fromStore, toStore)
			Expect(err).NotTo(HaveOccurred())
		})

		It("migrates data from fromStore to toStore", func() {
			Expect(toStore.CreateInstanceDetailsCallCount()).To(Equal(2))
			id1, serviceInstance1 := toStore.CreateInstanceDetailsArgsForCall(0)
			id2, serviceInstance2 := toStore.CreateInstanceDetailsArgsForCall(1)
			Expect([]string{id1, id2}).To(Equal([]string{"123", "456"}))
			Expect([]string{serviceInstance1.ServiceID, serviceInstance2.ServiceID}).To(Equal([]string{"some-service-1", "some-service-2"}))
		})
	})

	Context("when there are binding details in fromStore", func() {
		BeforeEach(func() {
			fromStore.RetrieveAllBindingDetailsReturns(map[string]brokerapi.BindDetails{
				"123": brokerapi.BindDetails{AppGUID: "some-app-1"},
				"456": brokerapi.BindDetails{AppGUID: "some-app-2"},
			}, nil)
			err := migrator.Migrate(fromStore, toStore)
			Expect(err).NotTo(HaveOccurred())
		})

		It("migrates data from fromStore to toStore", func() {
			Expect(toStore.CreateBindingDetailsCallCount()).To(Equal(2))
			id1, bindDetails1 := toStore.CreateBindingDetailsArgsForCall(0)
			id2, bindDetails2 := toStore.CreateBindingDetailsArgsForCall(1)
			Expect([]string{id1, id2}).To(Equal([]string{"123", "456"}))
			Expect([]string{bindDetails1.AppGUID, bindDetails2.AppGUID}).To(Equal([]string{"some-app-1", "some-app-2"}))
		})
	})
})
