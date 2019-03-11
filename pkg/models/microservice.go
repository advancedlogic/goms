package models

import (
	"fmt"
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/advancedlogic/goms/pkg/modules"
	"github.com/ankit-arora/go-utils/go-shutdown-hook"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Microservice struct {
	*modules.Environment
	mid       string
	transport interfaces.Transport
	discovery interfaces.Registry
	broker    interfaces.Broker
}

type MicroserviceBuilder struct {
	*Microservice
	modules.Exception
}

func NewMicroserviceBuilder(environment *modules.Environment) *MicroserviceBuilder {
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

func (mb *MicroserviceBuilder) WithTransport(transport interfaces.Transport) *MicroserviceBuilder {
	mb.transport = transport
	return mb
}

func (mb *MicroserviceBuilder) WithRestTransport() *MicroserviceBuilder {
	transport, err := modules.
		NewRestBuilder(mb.Environment).
		WithPort(mb.Environment.GetIntOrDefault("transport.port", 8080)).Build()
	if err != nil {
		mb.Fatal(err)
	}

	mb.transport = transport
	return mb
}

func (mb *MicroserviceBuilder) WithSubscription(topic string, callback func(interface{})) *MicroserviceBuilder {
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

func (mb *MicroserviceBuilder) Build() (*Microservice, error) {
	if err := mb.CheckErrors(mb.Errors()); err != nil {
		return nil, err
	}
	return mb.Microservice, nil
}

func (m *Microservice) Subscribe(topic string, callback func(interface{})) error {
	return m.broker.Subscribe(topic, callback)
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

func (m *Microservice) Transport() interfaces.Transport {
	return m.transport
}

func (m *Microservice) RestTransport() *modules.Rest {
	return m.transport.(*modules.Rest)
}

func (m *Microservice) Run() {
	if m.discovery != nil {
		err := m.discovery.Register()
		if err != nil {
			logrus.Error(err)
		}
	}

	if m.broker != nil {
		err := m.broker.Connect(m.GetStringOrDefault("broker.host", "localhost:4222"))
		if err != nil {
			logrus.Error(err)
		}
		err = m.broker.Subscribe(m.mid, func(i interface{}) {

		})
	}

	go_shutdown_hook.ADD(func() {
		m.Info("Goodbye and thanks for all the fish")
		err := m.transport.Stop()
		if err != nil {
			m.Fatal(err)
		}
	})
	err := m.transport.Run()
	if err != nil {
		m.Fatal(err)
	}
	go_shutdown_hook.Wait()
}
