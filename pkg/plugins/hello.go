package plugins

import (
	"errors"
	"fmt"
)

type Hello struct {
	Name string
}

func NewHello(name string) *Hello {
	return &Hello{
		Name: name,
	}
}

func (h *Hello) Process(descriptor interface{}) error {
	switch descriptor.(type) {
	case string:
		fmt.Printf("Hello %s -> %s", h.Name, descriptor.(string))
	default:
		return errors.New("argument must be a string")
	}
	return nil
}
