package config

import (
	"log"

	"github.com/kylelemons/go-gypsy/yaml"
)

var config *yaml.File

func init() {
	var err error
	config, err = yaml.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
}

func GetConf(key string, value string) (data string) {
	val, err := config.Get(key)

	if err != nil {
		val = value
	}

	data = val

	return data
}
