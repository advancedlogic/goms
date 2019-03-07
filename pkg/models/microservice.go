package models

import (
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/advancedlogic/goms/pkg/modules"
)

type Microservice struct {
	environment *modules.Environment
	transport   interfaces.Transport
}

func NewMicroservice(environment *modules.Environment) (*Microservice, error) {
	port := environment.GetIntOrDefault("transport.port", 8080)
	r, err := modules.NewRestBuilder(environment).WithPort(port).Build()
	if err != nil {
		return nil, err
	}
	return &Microservice{
		transport: r,
	}, nil
}
