package interfaces

type Transport interface {
	Run(...string) error
	Stop() error
	GetHandler(string, interface{})
	PostHandler(string, interface{})
	PutHandler(string, interface{})
	DeleteHandler(string, interface{})
}
