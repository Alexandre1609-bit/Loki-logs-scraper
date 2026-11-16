package main

import (
	"testing"
)

func resetMasker() {
	ipCounter = 0
	ipMap = make(map[string]string)
}

func TestMaskIpReal(t *testing.T) {
	resetMasker()
	got := maskIp("connexion depuis 10.0.0.1")
	want := "connexion depuis 192.0.2.1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaskIpInvalid(t *testing.T) {
	resetMasker()
	got := maskIp("999.999.999.999")
	want := "999.999.999.999"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaskIpString(t *testing.T) {
	resetMasker()
	got := maskIp("hello")
	want := "hello"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaskIpOther(t *testing.T) {
	resetMasker()
	got := maskIp("10.0.0.2 10.0.0.3")
	want := "192.0.2.1 192.0.2.2"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaskIpZeroNorm(t *testing.T) {
	resetMasker()
	got := maskIp("010.000.000.001 puis 10.0.0.1")
	want := "192.0.2.1 puis 192.0.2.1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
