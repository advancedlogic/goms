package models

import (
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/advancedlogic/goms/pkg/modules"
)

type Microservice struct {
	Configuration interfaces.Configuration
	Sync          interfaces.Sync
}

func NewMicroservice(environment *modules.Environment) (*Microservice, error) {
	r, err := modules.NewRestBuilder(environment).WithPort(9000).Build()
	if err != nil {
		return nil, err
	}
	return &Microservice{
		Sync: r,
	}, nil
}
