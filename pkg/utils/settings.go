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
	ServiceAccountName  string `default:"deploy-update-controller" split_words:"true"`
	ControllerNamespace string `default:"deploy-update-controller" split_words:"true"`
	InCluster           bool   `default:"false" split_words:"true"`
}

func NewSettings() *Specification {
	s := Specification{}
	// err := envconfig.Process("deployUpdater", &s)
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(s.InCluster)
	return &s
}
