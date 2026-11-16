package main

import "testing"

func TestMaskIpInvalid(t *testing.T) {
	got := maskIp("999.999.999.999")
	want := "999.999.999.999"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestMaskIp(t *testing.T) {
	got := maskIp("hello")
	want := "hello"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
