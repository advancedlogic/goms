package main

import (
	"github.com/advancedlogic/goms/examples/spider/plugin"
	"github.com/advancedlogic/goms/pkg/helpers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	builder := helpers.DefaultBuilder("spider", helpers.REGISTRY, helpers.TRANSPORT)

	if builder != nil {
		p := &plugin.Spider{}
		microservice, _ := builder.
			WithPlugin(p).
			Build()
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
		microservice.Run()
	} else {
		logrus.Fatal("Something wrong happened")
	}
}
