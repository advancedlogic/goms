package plugins

import (
	"errors"
	"fmt"
	"github.com/advancedlogic/goms/pkg/interfaces"
)

type Hello struct {
	interfaces.Service
	Name string
}

func NewHello(name string) *Hello {
	return &Hello{
		Name: name,
	}
}

func (r *Hello) Init(service interfaces.Service) error {
	r.Service = service
	return nil
}
func (r *Hello) Close() error { return nil }

func (h *Hello) Process(descriptor interface{}) (interface{}, error) {
	switch descriptor.(type) {
	case string:
		fmt.Printf("Hello %s -> %s", h.Name, descriptor.(string))
	default:
		return nil, errors.New("argument must be a string")
	}
	return "hello", nil
}
