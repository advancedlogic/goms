package modules

import (
	"errors"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"os"
)

type ConsulRegistry struct {
	id   string
	name string
	port int
}

type ConsulRegistryBuilder struct {
	*ConsulRegistry
	errors []error
}

func NewConsulRegistryBuilder() *ConsulRegistryBuilder {
	return &ConsulRegistryBuilder{
		ConsulRegistry: &ConsulRegistry{},
		errors:         make([]error, 0),
	}
}

func (crb *ConsulRegistryBuilder) WithID(id string) *ConsulRegistryBuilder {
	if id == "" {
		crb.errors = append(crb.errors, errors.New("id cannot be empty"))
	}
	crb.id = id
	return crb
}

func (crb *ConsulRegistryBuilder) WithName(name string) *ConsulRegistryBuilder {
	if name == "" {
		crb.errors = append(crb.errors, errors.New("name cannot be empty"))
	}
	crb.name = name
	return crb
}

func (crb *ConsulRegistryBuilder) WithHealthCheckingPort(port int) *ConsulRegistryBuilder {
	if port == 0 {
		crb.errors = append(crb.errors, errors.New("port must be greater than 0"))
	}
	crb.port = port
	return crb
}

func (crb *ConsulRegistryBuilder) Build() (*ConsulRegistry, error) {
	if len(crb.errors) > 0 {
		msg := ""
		for _, err := range crb.errors {
			msg = msg + fmt.Sprintf("%s\n", err.Error())
		}
		return nil, errors.New(msg)
	}
	return crb.ConsulRegistry, nil
}

func (cr *ConsulRegistry) Register() error {
	hostname := func() string {
		hn, err := os.Hostname()
		if err != nil {
			log.Fatalln(err)
		}
		return hn
	}

	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = cr.id
	registration.Name = cr.name
	address := hostname()
	registration.Address = address
	registration.Port = cr.port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck",
		address, cr.port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	return consul.Agent().ServiceRegister(registration)
}
