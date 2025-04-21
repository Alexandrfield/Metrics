package common

import (
	"testing"
)

func TestGetDataFromFile(t *testing.T) {
	d := GetDataFromFile("")
	if len(d) != 0 {
		t.Errorf("Expected 0 len")
	}
}
func TestGetDataFromBadFile(t *testing.T) {
	d := GetDataFromFile("test")
	if len(d) != 0 {
		t.Errorf("Expected 0 len")
	}
}
