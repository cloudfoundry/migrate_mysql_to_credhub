package migrator

import (
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/service-broker-store/brokerstore"
)

//go:generate counterfeiter -o fakes/fake_retirable_store.go . RetirableStore
type RetirableStore interface {
	Retire() error
	IsRetired() (bool, error)
	brokerstore.Store
}

//go:generate counterfeiter -o fakes/fake_activatable_store.go . ActivatibleStore
type ActivatableStore interface {
	Activate() error
	IsActivated() (bool, error)
	brokerstore.Store
}
type Migrator interface {
	Migrate(RetirableStore, ActivatableStore) error
}

type migrator struct {
	logger lager.Logger
}

func NewMigrator(logger lager.Logger) Migrator {
	return &migrator{
		logger: logger,
	}
}

func (m *migrator) Migrate(fromStore RetirableStore, toStore ActivatableStore) error {
	logger := m.logger.Session("migrate")
	logger.Info("start")
	defer logger.Info("end")

	activated, err := toStore.IsActivated()
	if err != nil {
		logger.Error("failed-to-check-if-credhub-is-activated", err)
		return err
	}

	if activated {
		logger.Info("credhub-already-activated")
		return nil
	}

	instanceDetails, err := fromStore.RetrieveAllInstanceDetails()
	if err != nil {
		logger.Error("failed-to-retrieve-all-instance-details", err)
		return err
	}

	logger.Info("instance-details", lager.Data{"count": len(instanceDetails)})
	for id, details := range instanceDetails {
		err = toStore.CreateInstanceDetails(id, details)
		if err != nil {
			logger.Error("failed-to-create-instance-details", err, lager.Data{"id": id, "service-details": details})
			return err
		}
	}
	bindingDetails, err := fromStore.RetrieveAllBindingDetails()
	if err != nil {
		logger.Error("failed-to-retrieve-all-binding-details", err)
		return err
	}

	logger.Info("binding-details", lager.Data{"count": len(bindingDetails)})
	for id, details := range bindingDetails {
		err = toStore.CreateBindingDetails(id, details)
		if err != nil {
			logger.Error("failed-to-create-binding-details", err, lager.Data{"id": id, "binding-details": details})
			return err
		}
	}

	err = toStore.Activate()
	if err != nil {
		return err
	}

	return fromStore.Retire()
}
