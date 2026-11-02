package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type conf struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

func (c *conf) loadConfig() *conf { //Se renseigner sur cette méthode de déclaration
	yamlFile, err := os.ReadFile("example-data.yaml")
	if err != nil {
		log.Fatalf("yamlFile.Get err #%v", err)
	}

	err = yaml.Unmarshal(yamlFile, c) //Faire la struct associée
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func mask() {
	c := &conf{}
	c.loadConfig()
	log.Print("Adresse: "+c.Address+":", c.Port)
}
