package modules

import (
	"errors"
	"github.com/nats-io/go-nats"
)

type Nats struct {
	endpoint            string
	conn                *nats.Conn
	userCredentialsPath string
	userJWT             string
	userNK              string
}

type NatsBuilder struct {
	*Environment
	*Nats
	errors []error
}

func NewNatsBuilder(environment *Environment) *NatsBuilder {
	return &NatsBuilder{
		Environment: environment,
		Nats:        &Nats{},
		errors:      make([]error, 0),
	}
}

func (nb *NatsBuilder) WithEndpoint(endpoint string) *NatsBuilder {
	if endpoint == "" {
		nb.errors = append(nb.errors, errors.New("endpoint cannot be empty"))
	}
	nb.endpoint = endpoint
	return nb
}

func (nb *NatsBuilder) WithUserCredentialsPath(path string) *NatsBuilder {
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

func (n *Nats) Subscribe(topic string, callback func(interface{})) error {
	return n.Subscribe(topic, callback)
}

func (n *Nats) Close() {
	n.conn.Close()
}
