package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type conf struct {
}

func (c *conf) loadConfig() *conf {

	yamlFile, err := os.ReadFile("example-data.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c) //Faire la struct yaml pour ref le c
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
