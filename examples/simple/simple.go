package main

import (
	"github.com/advancedlogic/goms/pkg/models"
	"github.com/advancedlogic/goms/pkg/modules"
	"github.com/advancedlogic/goms/pkg/plugins"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
)

func main() {
	environment, err := models.
		NewEnvironmentBuilder().
		WithConfigurationFile("config").
		WithName("simple").
		Build()
	if err != nil {
		log.Fatal(err)
	}

	d, err := modules.NewConsulRegistryBuilder(environment).Build()
	if err != nil {
		log.Fatal(err)
	}

	b, err := modules.NewNatsBuilder(environment).Build()
	if err != nil {
		log.Fatal(err)
	}

	r, err := modules.NewRestBuilder(environment).Build()
	if err != nil {
		log.Fatal(err)
	}

	s, err := modules.NewMinioBuilder(environment).Build()
	if err != nil {
		log.Fatal(err)
	}

	p := plugins.NewHello("hello")

	microservice, err := models.
		NewMicroserviceBuilder(environment).
		WithDiscovery(d).
		WithBroker(b).
		WithTransport(r).
		WithStore(s).
		WithPlugin(p).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	defer microservice.Close()

	microservice.GetHandler("/ping", func(ctx *gin.Context) {
		ctx.String(202, "pong")
	})

	microservice.GetHandler("/test/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		if err := microservice.Process(name); err != nil {
			ctx.String(400, err.Error())
		}
		ctx.String(202, name)
	})

	if err := microservice.SubscribeDefault(func(msg *nats.Msg) {
		if err := microservice.Process(string(msg.Data)); err != nil {

		}
	}); err != nil {
		log.Fatal(err)
	}
	microservice.Run()
}
