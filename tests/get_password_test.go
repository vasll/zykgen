package zykgen

import (
	"zykgen"
	"testing"
)

var data = map[string]string{
	"S090Y00000000": "UJ4NKUJ8KP",
}

func TestGetPassword(t *testing.T) {
	for serial, password := range data {
		genPassword := zykgen.GetPassword(serial, 10, zykgen.Cosmopolitan)
		if genPassword != password {
			t.Fatalf("got %s instead of %s expected wpa key", genPassword, password)
		}
	}
}
