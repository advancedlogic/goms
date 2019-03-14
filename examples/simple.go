package main

import (
	"errors"
	"github.com/advancedlogic/goms/pkg/models"
	"github.com/advancedlogic/goms/pkg/modules"
	"github.com/advancedlogic/goms/pkg/plugins"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
)

func main() {
	environment, err := modules.NewEnvironmentBuilder().WithConfigurationFile("config").WithName("simple").Build()
	if err != nil {
		log.Fatal(err)
	}

	r, err := modules.
		NewRestBuilder(environment).
		WithPort(environment.GetIntOrDefault("transport.port", 8080)).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	d, err := modules.NewConsulRegistryBuilder(environment).
		WithName(environment.GetStringOrDefault("service.name", "default")).
		WithID(environment.GetStringOrDefault("service.id", "default")).
		WithHealthCheckingPort(environment.GetIntOrDefault("transport.port", 8080)).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	b, err := modules.NewNatsBuilder(environment).
		WithEndpoint("localhost:4222").Build()
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	microservice, err := models.
		NewMicroserviceBuilder(environment).
		WithTransport(r).
		WithDiscovery(d).
		WithBroker(b).
		WithPlugin(plugins.NewHello("hello")).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	microservice.GetHandler("/ping", func(ctx *gin.Context) {
		ctx.String(202, "pong")
	})

	microservice.GetHandler("/test/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		if name == "" {
			ctx.String(400, errors.New("param name cannot be empty").Error())
			return
		}
		if err := microservice.Process(name); err != nil {
			ctx.String(400, err.Error())
		}
		ctx.String(202, name)
	})

	if err := microservice.Subscribe("test", func(msg *nats.Msg) {

	}); err != nil {
		log.Fatal(err)
	}
	microservice.Run()
}
