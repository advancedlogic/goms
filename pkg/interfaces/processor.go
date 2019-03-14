package interfaces

type Processor interface {
	Process(interface{}) error
}
