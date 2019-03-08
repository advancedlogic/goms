package modules

import (
	"errors"
	"github.com/advancedlogic/goms/pkg/interfaces"
)

type Discovery struct {
	registry interfaces.Registry
}

type DiscoveryBuilder struct {
	*Environment
	*Discovery
	errors []error
}

func NewDiscoveryBuilder(environment *Environment) *DiscoveryBuilder {
	return &DiscoveryBuilder{
		Discovery:   &Discovery{},
		Environment: environment,
		errors:      make([]error, 0),
	}
}

func (db *DiscoveryBuilder) WithRegistry(registry interfaces.Registry) *DiscoveryBuilder {
	if registry == nil {
		db.errors = append(db.errors, errors.New("a registry must be specified"))
	}
	db.registry = registry
	return db
}

func (db *DiscoveryBuilder) Build() (*Discovery, error) {
	err := db.CheckErrors(db.errors)
	if err != nil {
		return nil, err
	}
	return db.Discovery, nil
}

func (d *Discovery) Register() error {
	return d.registry.Register()
}
