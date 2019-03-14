package interfaces

type Cache interface {
	Put(string, interface{}) error
	Get(string) (interface{}, error)
}
