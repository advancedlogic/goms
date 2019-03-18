package main

import (
	"github.com/advancedlogic/goms/pkg/helpers"
	"github.com/sirupsen/logrus"
)

func main() {
	builder := helpers.DefaultMicroserviceBuilder("frontend", helpers.TRANSPORT, helpers.REGISTRY)

	if builder != nil {
		microservice, err := builder.Build()
		if err != nil {
			logrus.Fatal(err)
		}
		microservice.Run()
	}
}
