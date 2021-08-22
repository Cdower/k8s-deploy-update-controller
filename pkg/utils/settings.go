package utils

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Channel             string `required:"true" default:"https://github.com/home-assistant/core/releases/latest"`
	Namespace           string `required:"true" default:"default"`
	Deployment          string `required:"true" default:"home-assistant"`
	Container           string `required:"true" default:"main"`
	ServiceAccountName  string `default:"deploy-update-controller"`
	ControllerNamespace string `default:"deploy-update-controller"`
	InCluster           bool   `default:"false"`
}

func NewSettings() *Specification {
	s := Specification{}
	err := envconfig.Process("deployUpdater", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &s
}
