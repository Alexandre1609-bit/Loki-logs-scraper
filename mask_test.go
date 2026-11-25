package main

import (
	"testing"
)

func resetMasker() {
	ipCounter = 0
	ipMap = make(map[string]string)
}

func TestMaskIp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "BasicIp", input: "10.0.0.1", expected: "192.0.2.1"},
		{name: "InvalidIp", input: "999.999.999.999", expected: "999.999.999.999"},
		{name: "String", input: "Hello", expected: "Hello"},
		{name: "MultipleIp", input: "10.0.0.2 10.0.0.3", expected: "192.0.2.1 192.0.2.2"},
		{name: "ZeroNormSameIp", input: "010.000.000.001 and 10.0.0.1", expected: "192.0.2.1 and 192.0.2.1"},
		{name: "TextAroundIp", input: "connexion from 10.0.0.1", expected: "connexion from 192.0.2.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetMasker()
			result := maskIp(tt.input)
			if result != tt.expected {
				t.Errorf("maskIp(%q) = %q; attendu %q", tt.input, result, tt.expected)
			}
		})
	}
}
