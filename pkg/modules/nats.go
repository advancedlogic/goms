package modules

import (
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
			logger: environment.Logger,
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

func (n *Nats) Connect(endpoint string) error {
	var err error
	if conn, err := nats.Connect(endpoint); err == nil {
		n.conn = conn
		return nil
	}
	return err

}

func (n *Nats) Publish(topic string, message []byte) error {
	return n.conn.Publish(topic, message)
}

func (n *Nats) Subscribe(topic string, callback func(interface{})) error {
	return n.Subscribe(topic, callback)
}

func (n *Nats) Close() {
	n.conn.Close()
}
