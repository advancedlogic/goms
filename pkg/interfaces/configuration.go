package interfaces

type Configuration interface {
	Load(path string)
	Save(path string)
}
