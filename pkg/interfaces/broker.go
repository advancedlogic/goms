package interfaces

import "github.com/nats-io/go-nats"

type Broker interface {
	Run() error
	Endpoint() string
	Connect() error
	Publish(string, []byte) error
	Subscribe(string, nats.MsgHandler) error
	Unsubscribe(string) error
	Close()
}
