package modules

import (
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
	*Environment
	*ConsulRegistry
	Exception
}

func NewConsulRegistryBuilder(environment *Environment) *ConsulRegistryBuilder {
	return &ConsulRegistryBuilder{
		ConsulRegistry: &ConsulRegistry{},
	}
}

func (crb *ConsulRegistryBuilder) WithID(id string) *ConsulRegistryBuilder {
	if id == "" {
		crb.Catch("id cannot be empty")
	}
	crb.id = id
	return crb
}

func (crb *ConsulRegistryBuilder) WithName(name string) *ConsulRegistryBuilder {
	if name == "" {
		crb.Catch("name cannot be empty")
	}
	crb.ConsulRegistry.name = name
	return crb
}

func (crb *ConsulRegistryBuilder) WithHealthCheckingPort(port int) *ConsulRegistryBuilder {
	if port == 0 {
		crb.Catch("port must be greater than 0")
	}
	crb.port = port
	return crb
}

func (crb *ConsulRegistryBuilder) Build() (*ConsulRegistry, error) {
	if err := crb.CheckErrors(crb.errors); err != nil {
		return nil, err
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
