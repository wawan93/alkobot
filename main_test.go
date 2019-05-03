package main

import "testing"

func TestGetRandomPhrase(t *testing.T) {
	got := GetRandomPhrase()
	if got == "" {
		t.Error("empty string returned")
	}
}
