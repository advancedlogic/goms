package interfaces

type Microservice interface {
	BrokerHandler(string, func(interface{}) error)
	AddRestHandler(string, func(interface{}) error)
}
