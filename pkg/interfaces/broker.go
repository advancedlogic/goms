package interfaces

type Broker interface {
	Connect(string) error
	Publish(string, []byte) error
	Subscribe(string, func(interface{})) error
	Close()
}
