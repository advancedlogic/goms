package main

import (
	"github.com/advancedlogic/goms/pkg/models"
	"github.com/advancedlogic/goms/pkg/modules"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	environment, err := modules.NewEnvironmentBuilder().WithConfigurationFile("config").WithName("simple").Build()
	if err != nil {
		log.Fatal(err)
	}

	/*r, err := modules.
		NewRestBuilder(environment).
		WithPort(environment.GetIntOrDefault("transport.port", 8080)).
		WithGetHandler("/ping", func(ctx *gin.Context) {
			ctx.String(202, "pong")
		}).
		Build()
	if err != nil {
		log.Fatal(err)
	}*/

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

	microservice, err := models.
		NewMicroserviceBuilder(environment).
		WithRestTransport().
		WithDiscovery(d).
		WithBroker(b).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	microservice.RestTransport().GetHandler("ping", func(ctx *gin.Context) {

	})
	microservice.Run()
}
