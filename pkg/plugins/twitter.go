package plugins

import "github.com/advancedlogic/goms/pkg/interfaces"

type TwitterSource struct {
}

type Twitter struct {
	interfaces.Service
}

type Source struct {
}

func (t *Twitter) Init(service interfaces.Service) error {
	t.Service = service
	return nil
}

func (t *Twitter) Close() error {
	return nil
}

func (t *Twitter) Process(source interface{}) (interface{}, error) {
	return nil, nil
}
