package main

import (
	"testing"
)

func TestTest(t *testing.T) {
	got := 1
	want := 1

	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
}

func TestFailing(t *testing.T) {
	got := pisscum(4)
	want := 4

	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
}

func TestTest3(t *testing.T) {
	got := pisscum(4)
	want := 8

	if got != want {
		t.Errorf("Got %d, wanted %d", got, want)
	}
}
