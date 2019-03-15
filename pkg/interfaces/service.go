package interfaces

type Service interface {
	//Modules
	Transport() Transport
	Registry() Registry
	Broker() Broker
	Store() Store
	Cache() Cache

	Running() bool

	//Configuration
	Config(path string, defaultValue interface{}) interface{}

	//Log
	Error(string)
	Errorf(string, ...interface{})
	Info(string)
	Infof(string, ...interface{})
	Warn(string)
	Warnf(string, ...interface{})
	Fatal(string)
	Fatalf(string, ...interface{})
	Debug(string)
	Debugf(string, ...interface{})
}
