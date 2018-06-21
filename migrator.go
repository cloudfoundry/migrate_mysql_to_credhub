package main

import (
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/service-broker-store/brokerstore"
)

type Migrator interface {
	Migrate(brokerstore.Store, brokerstore.Store) error
}

type migrator struct {
	logger lager.Logger
}

func NewMigrator(logger lager.Logger) Migrator {
	return &migrator{
		logger: logger,
	}
}

func (m *migrator) Migrate(fromStore brokerstore.Store, toStore brokerstore.Store) error {
	instanceDetails, err := fromStore.RetrieveAllInstanceDetails()
	if err != nil {
		m.logger.Error("failed-to-retrieve-all-instance-details", err)
		return err
	}
	for id, details := range instanceDetails {
		err = toStore.CreateInstanceDetails(id, details)
		if err != nil {
			m.logger.Error("failed-to-create-instance-details", err, lager.Data{"id": id, "service-details": details})
			return err
		}
	}
	bindingDetails, err := fromStore.RetrieveAllBindingDetails()
	if err != nil {
		m.logger.Error("failed-to-retrieve-all-binding-details", err)
		return err
	}
	for id, details := range bindingDetails {
		err = toStore.CreateBindingDetails(id, details)
		if err != nil {
			m.logger.Error("failed-to-create-binding-details", err, lager.Data{"id": id, "binding-details": details})
			return err
		}
	}
	return nil
}
