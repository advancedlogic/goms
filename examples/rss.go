package main

import (
	"github.com/advancedlogic/goms/pkg/models"
	"github.com/advancedlogic/goms/pkg/modules"
	"github.com/advancedlogic/goms/pkg/plugins"
	"github.com/advancedlogic/goms/pkg/tools"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/go-nats"
	"log"
	"net/http"
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

	p := plugins.NewRSS()

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

	microservice.StaticFilesFolder("/static", "./www")
	microservice.PostHandler("/api/v1/rss/source", func(ctx *gin.Context) {
		var descriptor plugins.Descriptor
		if ctx.ShouldBind(&descriptor) == nil {
			feeds, err := microservice.Process(descriptor)
			if err != nil {
				ctx.String(http.StatusBadGateway, err.Error())
				return
			}
			for _, feed := range feeds.([]string) {
				id := tools.SHA1(feed)
				microservice.Infof("[%s] %s", id, feed)
				if err := microservice.Create(id, feed); err != nil {
					ctx.String(http.StatusBadGateway, err.Error())
					return
				}
			}
			ctx.JSON(http.StatusOK, descriptor)
		}
	})

	if err := microservice.SubscribeDefault(func(msg *nats.Msg) {
		feeds, err := microservice.Process(string(msg.Data))
		if err != nil {
			microservice.Error(err.Error())
			return
		}
		for _, feed := range feeds.([]string) {
			id := tools.SHA1(feed)
			microservice.Infof("[%s] %s", id, feed)
			if err := microservice.Create(id, feed); err != nil {
				microservice.Error(err.Error())
			}
		}
	}); err != nil {
		log.Fatal(err)
	}
	microservice.Run()
}
