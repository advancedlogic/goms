package models

import (
	"fmt"
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/ankit-arora/go-utils/go-shutdown-hook"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nats-io/go-nats"
	"time"
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
	running   bool
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

func (mb *MicroserviceBuilder) WithTransport(transport interfaces.Transport) *MicroserviceBuilder {
	mb.transport = transport
	return mb
}

func (mb *MicroserviceBuilder) WithSubscription(topic string, callback nats.MsgHandler) *MicroserviceBuilder {
	err := mb.broker.Subscribe(topic, callback)
	if err != nil {
		mb.Fatal(err.Error())
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
	if err := mb.CheckErrors(mb.Exception.Errors()); err != nil {
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

func (m *Microservice) Middleware(middleware func(*gin.Context)) {
	m.transport.Middleware(middleware)
}

func (m *Microservice) StaticFilesFolder(uri, folder string) {
	m.transport.StaticFilesFolder(uri, folder)
}

func (m *Microservice) Publish(topic string, handler interface{}) error {
	return m.broker.Publish(topic, handler)
}

func (m *Microservice) Process(instance interface{}) (interface{}, error) {
	return m.module.Process(instance)
}

func (m *Microservice) Create(key string, value interface{}) error {
	return m.store.Create(key, value)
}

func (m *Microservice) Delete(key string) error {
	return m.store.Delete(key)
}

func (m *Microservice) Update(key string, value interface{}) error {
	return m.store.Update(key, value)
}

func (m *Microservice) Read(key string) (interface{}, error) {
	return m.store.Read(key)
}

func (m *Microservice) List() ([]interface{}, error) {
	return m.store.List()
}

func (m *Microservice) Errors() []error {
	return m.Errors()
}

func (m *Microservice) Fatal(msg string) {
	m.Logger.Error(msg)
}

func (m *Microservice) Fatalf(msg string, args ...interface{}) {
	m.Logger.Fatalf(msg, args)
}

func (m *Microservice) Error(msg string) {
	m.Logger.Error(msg)
}

func (m *Microservice) Errorf(msg string, args ...interface{}) {
	m.Logger.Errorf(msg, args)
}

func (m *Microservice) Warn(msg string) {
	m.Logger.Warn(msg)
}

func (m *Microservice) Warnf(msg string, args ...interface{}) {
	m.Logger.Warnf(msg, args)
}

func (m *Microservice) Info(msg string) {
	m.Logger.Info(msg)
}
func (m *Microservice) Infof(msg string, args ...interface{}) {
	m.Logger.Infof(msg, args)
}

func (m *Microservice) Debug(msg string) {
	m.Logger.Debug(msg)
}

func (m *Microservice) Debugf(msg string, args ...interface{}) {
	m.Logger.Debugf(msg, args)
}

func (m *Microservice) Config(path string, defaultValue interface{}) interface{} {
	switch defaultValue.(type) {
	case string:
		return m.GetStringOrDefault(path, defaultValue.(string))
	case int:
		return m.GetIntOrDefault(path, defaultValue.(int))
	case bool:
		return m.GetBoolOrDefault(path, defaultValue.(bool))
	case time.Duration:
		return m.GetDurationOrDefault(path, defaultValue.(time.Duration))
	}

	return nil
}

func (mb *MicroserviceBuilder) WithPlugin(processor interfaces.Processor) *MicroserviceBuilder {
	mb.module = processor
	return mb
}

func (m *Microservice) Run() {
	defer m.Close()

	go_shutdown_hook.ADD(func() {
		m.Info("Goodbye and thanks for all the fish")
		err := m.transport.Stop()
		if err != nil {
			m.Fatal(err.Error())
		}
	})

	if m.discovery != nil {
		if err := m.discovery.Register(); err != nil {
			m.Error(err.Error())
		}
		m.Info("Discovery initialized")
	}

	if m.broker != nil {
		if err := m.broker.Run(); err != nil {
			m.Error(err.Error())
		}
		m.Info("Broker initialized")
	}
	m.running = true

	if m.module != nil {
		if err := m.module.Init(m); err != nil {
			m.Error(err.Error())
		}
		m.Info("Module initialized")
	}

	if m.transport != nil {
		if err := m.transport.Run(); err != nil {
			m.Error(err.Error())
		}
		m.Info("Transport initialized")
	}

	go_shutdown_hook.Wait()
}

func (m *Microservice) Running() bool {
	return m.running
}

func (m *Microservice) Close() {
	m.running = false
	if m.broker != nil {
		m.broker.Close()
	}
	if m.transport != nil {
		if err := m.transport.Stop(); err != nil {
			m.Fatal(err.Error())
		}
	}
	if m.module != nil {
		if err := m.module.Close(); err != nil {
			m.Fatal(err.Error())
		}
	}
}

func (m *Microservice) Transport() interfaces.Transport {
	return m.transport
}

func (m *Microservice) Broker() interfaces.Broker {
	return m.broker
}

func (m *Microservice) Registry() interfaces.Registry {
	return m.discovery
}

func (m *Microservice) Cache() interfaces.Cache {
	return m.cache
}

func (m *Microservice) Store() interfaces.Store {
	return m.store
}
