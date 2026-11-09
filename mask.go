package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var findIp = regexp.MustCompile(`\b\d+(?:\.\d+){3,}\b`)
var ipCounter = 0
var ipMap = make(map[string]string)

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

	rst := findIp.ReplaceAllStringFunc(addresse, func(match string) string {
		octets := strings.Split(match, ".")
		for index := range octets {
			octets[index] = strings.TrimLeft(octets[index], "0")
			if octets[index] == "" {
				octets[index] = "0"
			}
		}
		joinIp := strings.Join(octets, ".")
		ip := net.ParseIP(joinIp)
		if ip == nil {
			return match
		}
		if value, ok := ipMap[joinIp]; ok {
			return value
		} else {
			ipCounter++
			alias := fmt.Sprintf("192.0.2.%d", ipCounter)
			ipMap[joinIp] = alias
			return alias
		}
	})
	return rst
}

func mask() {
	c := &conf{}
	c.loadConfig()
	log.Print("Adresse: "+c.Address+":", c.Port)
}
