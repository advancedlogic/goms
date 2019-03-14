package models

import (
	"fmt"
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/advancedlogic/goms/pkg/modules"
	"github.com/ankit-arora/go-utils/go-shutdown-hook"
	"github.com/google/uuid"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"log"
)

type Microservice struct {
	*Environment
	mid       string
	transport interfaces.Transport
	discovery interfaces.Registry
	broker    interfaces.Broker
	module    interfaces.Processor
	store     interfaces.Store
	cache     interfaces.Cache
}

type MicroserviceBuilder struct {
	*Microservice
	Exception
}

func NewMicroserviceBuilder(environment *Environment) *MicroserviceBuilder {
	mid, err := uuid.NewUUID()
	if err != nil {
		environment.Fatal(err)
	}
	return &MicroserviceBuilder{
		Microservice: &Microservice{
			mid:         fmt.Sprintf("%s-$s", environment.GetStringOrDefault("service.name", "default"), mid.String()),
			Environment: environment,
		},
	}
}

func (mb *MicroserviceBuilder) Default() *MicroserviceBuilder {
	r, err := modules.NewRestBuilder(mb.Environment).Build()
	if err != nil {
		log.Fatal(err)
	}

	d, err := modules.NewConsulRegistryBuilder(mb.Environment).Build()
	if err != nil {
		log.Fatal(err)
	}

	b, err := modules.NewNatsBuilder(mb.Environment).Build()
	if err != nil {
		log.Fatal(err)
	}

	return mb.WithTransport(r).WithDiscovery(d).WithBroker(b)
}

func (mb *MicroserviceBuilder) WithTransport(transport interfaces.Transport) *MicroserviceBuilder {
	mb.transport = transport
	return mb
}

func (mb *MicroserviceBuilder) WithSubscription(topic string, callback nats.MsgHandler) *MicroserviceBuilder {
	err := mb.broker.Subscribe(topic, callback)
	if err != nil {
		mb.Fatal(err)
	}
	return mb
}

func (mb *MicroserviceBuilder) WithDiscovery(discovery interfaces.Registry) *MicroserviceBuilder {
	mb.discovery = discovery
	return mb
}

func (mb *MicroserviceBuilder) WithBroker(broker interfaces.Broker) *MicroserviceBuilder {
	mb.broker = broker
	return mb
}

func (mb *MicroserviceBuilder) WithStore(store interfaces.Store) *MicroserviceBuilder {
	mb.store = store
	return mb
}

func (mb *MicroserviceBuilder) WithCache(cache interfaces.Cache) *MicroserviceBuilder {
	mb.cache = cache
	return mb
}

func (mb *MicroserviceBuilder) Build() (*Microservice, error) {
	if err := mb.CheckErrors(mb.Errors()); err != nil {
		return nil, err
	}
	return mb.Microservice, nil
}

func (m *Microservice) Subscribe(topic string, callback nats.MsgHandler) error {
	return m.broker.Subscribe(topic, callback)
}

func (m *Microservice) SubscribeDefault(handler nats.MsgHandler) error {
	return m.broker.Subscribe(m.mid, handler)
}

func (m *Microservice) GetHandler(endpoint string, handler interface{}) {
	m.transport.GetHandler(endpoint, handler)
}

func (m *Microservice) PostHandler(endpoint string, handler interface{}) {
	m.transport.PostHandler(endpoint, handler)
}

func (m *Microservice) PutHandler(endpoint string, handler interface{}) {
	m.transport.PutHandler(endpoint, handler)
}

func (m *Microservice) DeleteHandler(endpoint string, handler interface{}) {
	m.transport.DeleteHandler(endpoint, handler)
}

func (m *Microservice) Process(instance interface{}) error {
	return m.module.Process(instance)
}

func (m *Microservice) Transport() interfaces.Transport {
	return m.transport
}

func (mb *MicroserviceBuilder) WithPlugin(processor interfaces.Processor) *MicroserviceBuilder {
	mb.module = processor
	return mb
}

func (m *Microservice) Run() {
	if m.discovery != nil {
		err := m.discovery.Register()
		if err != nil {
			logrus.Error(err)
		}
	}

	if m.broker != nil {
		err := m.broker.Run()
		if err != nil {
			logrus.Error(err)
		}
	}

	err := m.transport.Run()
	if err != nil {
		m.Fatal(err)
	}

	go_shutdown_hook.ADD(func() {
		m.Info("Goodbye and thanks for all the fish")
		err := m.transport.Stop()
		if err != nil {
			m.Fatal(err)
		}
	})
	go_shutdown_hook.Wait()
}

func (m *Microservice) Close() {
	if m.broker != nil {
		m.broker.Close()
	}
	if m.transport != nil {
		err := m.transport.Stop()
		if err != nil {
			m.Logger.Fatal(err)
		}
	}
}
