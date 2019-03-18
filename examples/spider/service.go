package main

import (
	"encoding/json"
	"github.com/advancedlogic/goms/examples/spider/plugin"
	"github.com/advancedlogic/goms/pkg/helpers"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Fatal(err)
		}
	}()

	builder := helpers.DefaultMicroserviceBuilder("spider", helpers.REGISTRY, helpers.TRANSPORT, helpers.BROKER)

	if builder != nil {
		p := &plugin.Spider{}
		microservice, _ := builder.
			WithPlugin(p).
			Build()

		//REST API Handler
		microservice.PostHandler("/api/v1/spider/url", func(ctx *gin.Context) {
			var source plugin.SpiderSource
			err := ctx.BindJSON(&source)
			if err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{
					"error": err.Error(),
				})
				return
			}
			descriptor, err := microservice.Process(source)
			if err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{
					"error": err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusOK, descriptor)
		})

		inTopic := microservice.GetStringOrDefault("broker.in", "spider")
		outTopic := microservice.GetStringOrDefault("broker.out", "langdec")

		//Broker topic Handler
		if err := microservice.Subscribe(inTopic, func(msg *nats.Msg) {
			var source plugin.SpiderDescriptor
			if err := json.Unmarshal(msg.Data, &source); err != nil {
				microservice.Error(err.Error())
			}
			descriptor, err := microservice.Process(source)
			if err != nil {
				microservice.Error(err.Error())
			}
			jsonDescriptor, err := json.Marshal(descriptor)
			if err != nil {
				microservice.Error(err.Error())
			}
			if err := microservice.Publish(outTopic, jsonDescriptor); err != nil {
				microservice.Error(err.Error())
			}
		}); err != nil {
			microservice.Fatal(err.Error())
		}

		microservice.Run()
	} else {
		logrus.Fatal("Something wrong happened")
	}
}
