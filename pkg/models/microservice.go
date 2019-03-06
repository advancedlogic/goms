package models

import (
	"github.com/advancedlogic/goms/pkg/interfaces"
	"github.com/advancedlogic/goms/pkg/rest"
)

type Microservice struct {
	Configuration interfaces.Configuration
	Sync          interfaces.Sync
}

func NewMicroservice() (*Microservice, error) {
	r, err := rest.NewBuilder().Port(900).Build()
	if err != nil {
		return nil, err
	}
	return &Microservice{
		Sync: r,
	}, nil
}
