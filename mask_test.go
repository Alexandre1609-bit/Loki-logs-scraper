package main

import (
	"testing"
)

func resetMasker() {
	ipCounter = 0
	ipMap = make(map[string]string)
}

func TestMaskIP(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "BasicIP", input: "10.0.0.1", expected: "192.0.2.1"},
		{name: "InvalidIP", input: "999.999.999.999", expected: "999.999.999.999"},
		{name: "String", input: "Hello", expected: "Hello"},
		{name: "MultipleIP", input: "10.0.0.2 10.0.0.3", expected: "192.0.2.1 192.0.2.2"},
		{name: "ZeroNormSameIP", input: "010.000.000.001 and 10.0.0.1", expected: "192.0.2.1 and 192.0.2.1"},
		{name: "TextAroundIP", input: "connexion from 10.0.0.1", expected: "connexion from 192.0.2.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetMasker()
			got := maskIP(tt.input)
			if got != tt.expected {
				t.Errorf("maskIP(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
