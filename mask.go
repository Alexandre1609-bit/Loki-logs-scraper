package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var findIP = regexp.MustCompile(`\b\d+(?:\.\d+){3,}\b`)
var ipCounter = 0
var ipMap = make(map[string]string)

func maskIP(line string) string {
	rst := findIP.ReplaceAllStringFunc(line, func(match string) string {
		octets := strings.Split(match, ".")
		for index := range octets {
			octets[index] = strings.TrimLeft(octets[index], "0")
			if octets[index] == "" {
				octets[index] = "0"
			}
		}
		joinIP := strings.Join(octets, ".")
		ip := net.ParseIP(joinIP)
		if ip == nil {
			return match
		}
		if value, ok := ipMap[joinIP]; ok {
			return value
		}
		ipCounter++
		alias := fmt.Sprintf("192.0.2.%d", ipCounter)
		ipMap[joinIP] = alias
		return alias
	})
	return rst
}
