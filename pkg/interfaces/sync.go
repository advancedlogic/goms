package interfaces

type Sync interface {
	Run(...string) error
}
