package common

import "testing"

func TestWarnf(t *testing.T) {
	fl := FakeLogger{}
	fl.Warnf("test:%d,%s,%t", 4, "zxc", true)
}
func TestInfof(t *testing.T) {
	fl := FakeLogger{}
	fl.Infof("test:%d,%s,%t", 4, "zxc", true)
}
func TestFatalf(t *testing.T) {
	fl := FakeLogger{}
	fl.Fatalf("test:%d,%s,%t", 4, "zxc", true)
}
