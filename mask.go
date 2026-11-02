package main

import (
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

var findIp = regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)

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

func maskIp(addresse string) string {

	return findIp.ReplaceAllString(addresse, "IP_MASQUEE")
}

func mask() {
	c := &conf{}
	c.loadConfig()
	log.Print("Adresse: "+c.Address+":", c.Port)
}
