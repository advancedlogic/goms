package helpers

import (
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/advancedlogic/goms/pkg/models"
	"github.com/advancedlogic/goms/pkg/modules"
)

const (
	TRANSPORT = "transport"
	BROKER    = "broker"
	REGISTRY  = "registry"
	STORE     = "store"
)

func DefaultTransport(builder *models.MicroserviceBuilder) error {
	r, err := modules.NewRestBuilder(builder.Environment).Build()
	if err != nil {
		return err
	}
	builder.WithTransport(r)
	return nil
}

func DefaultBroker(builder *models.MicroserviceBuilder) error {
	b, err := modules.NewNatsBuilder(builder.Environment).Build()
	if err != nil {
		return err
	}
	builder.WithBroker(b)
	return nil
}

func DefaultRegistry(builder *models.MicroserviceBuilder) error {
	c, err := modules.NewConsulRegistryBuilder(builder.Environment).Build()
	if err != nil {
		return err
	}
	builder.WithDiscovery(c)
	return nil
}

func DefaultStore(builder *models.MicroserviceBuilder) error {
	s, err := modules.NewMinioBuilder(builder.Environment).Build()
	if err != nil {
		return err
	}
	builder.WithStore(s)
	return nil
}

func Default(builder *models.MicroserviceBuilder) error {
	if err := DefaultBroker(builder); err != nil {
		return err
	}

	if err := DefaultRegistry(builder); err != nil {
		return err
	}

	if err := DefaultTransport(builder); err != nil {
		return err
	}

	if err := DefaultStore(builder); err != nil {
		return err
	}

	return nil
}

func DefaultN(builder *models.MicroserviceBuilder, modules ...interface{}) {
	for _, module := range modules {
		switch module.(type) {
		case interfaces.Transport:
			builder.WithTransport(module.(interfaces.Transport))
		case interfaces.Broker:
			builder.WithBroker(module.(interfaces.Broker))
		case interfaces.Registry:
			builder.WithDiscovery(module.(interfaces.Registry))
		case interfaces.Store:
			builder.WithStore(module.(interfaces.Store))
		}
	}
}

func DefaultS(builder *models.MicroserviceBuilder, modules ...string) {
	for _, module := range modules {
		switch module {
		case TRANSPORT:
			DefaultTransport(builder)
		case BROKER:
			DefaultBroker(builder)
		case REGISTRY:
			DefaultRegistry(builder)
		case STORE:
			DefaultStore(builder)
		}
	}
}

func DefaultMicroserviceBuilder(name string, modules ...string) *models.MicroserviceBuilder {
	environment, err := models.
		NewEnvironmentBuilder().
		WithConfigurationFile("config").
		WithName(name).
		Build()
	if err != nil {
		return nil
	}

	builder := models.NewMicroserviceBuilder(environment)
	DefaultS(builder, modules...)
	return builder
}
