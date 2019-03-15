package interfaces

type Transport interface {
	Run() error
	Stop() error
	GetHandler(string, interface{})
	PostHandler(string, interface{})
	PutHandler(string, interface{})
	DeleteHandler(string, interface{})
	Middleware(interface{})
	StaticFilesFolder(string, string)
}
