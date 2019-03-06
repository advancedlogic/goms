package interfaces

type Transport interface {
	Run(...string) error
}
