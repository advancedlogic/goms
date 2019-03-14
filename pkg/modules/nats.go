package modules

import (
	"errors"
	"fmt"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

type Nats struct {
	endpoint            string
	conn                *nats.Conn
	userCredentialsPath string
	userJWT             string
	userNK              string
	logger              *logrus.Logger
	handlers            map[string]nats.MsgHandler
	subscriptions       map[string]*nats.Subscription
}

type NatsBuilder struct {
	*Environment
	*Nats
	Exception
}

func NewNatsBuilder(environment *Environment) *NatsBuilder {
	return &NatsBuilder{
		Environment: environment,
		Nats: &Nats{
			logger:        environment.Logger,
			handlers:      make(map[string]nats.MsgHandler),
			subscriptions: make(map[string]*nats.Subscription),
		},
	}
}

func (nb *NatsBuilder) WithEndpoint(endpoint string) *NatsBuilder {
	if endpoint == "" {
		nb.Catch("endpoint cannot be empty")
	}
	nb.endpoint = endpoint
	return nb
}

func (nb *NatsBuilder) WithUserCredentialsPath(path string) *NatsBuilder {
	if path == "" {
		nb.Catch("user credentials path cannot be empty")
	}
	nb.userCredentialsPath = path
	return nb
}

func (nb *NatsBuilder) Build() (*Nats, error) {
	err := nb.CheckErrors(nb.errors)
	if err != nil {
		return nil, err
	}
	return nb.Nats, nil
}

func (n *Nats) Connect() error {
	var err error
	if conn, err := nats.Connect(n.endpoint); err == nil {
		n.conn = conn
		return nil
	}
	return err

}

func (n *Nats) Publish(topic string, message []byte) error {
	return n.conn.Publish(topic, message)
}

func (n *Nats) Subscribe(topic string, handler nats.MsgHandler) error {
	n.handlers[topic] = handler
	return nil
}

func (n *Nats) Unsubscribe(topic string) error {
	if subscription, exists := n.subscriptions[topic]; exists {
		return subscription.Unsubscribe()
	}
	return errors.New(fmt.Sprintf("topic %s does not exist", topic))
}

func (n *Nats) Run() error {
	if err := n.Connect(); err != nil {
		return err
	}
	for topic, handler := range n.handlers {
		subscription, err := n.conn.Subscribe(topic, handler)
		if err != nil {
			return err
		}
		n.subscriptions[topic] = subscription
	}
	return nil
}

func (n *Nats) Endpoint() string {
	return n.endpoint
}

func (n *Nats) Close() {
	n.conn.Close()
}
