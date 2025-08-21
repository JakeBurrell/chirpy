package main

import (
	"testing"
)

func TestResplaceProfanity(t *testing.T) {
	testText := "This is a kerfuffle opinion I need to share with the world"
	got := replaceProfanity(testText)
	want := "This is a **** opinion I need to share with the world"
	if got != want {
		t.Errorf("replaceProfanity(%s) = %s\n want = %s", testText, got, want)
	}
}
