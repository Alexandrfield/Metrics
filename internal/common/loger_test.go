package common

import "testing"

func TestLogWarnf(t *testing.T) {
	fl := FakeLogger{}
	fl.Warnf("test:%d,%s,%t", 4, "zxc", true)
}
func TestLogInfof(t *testing.T) {
	fl := FakeLogger{}
	fl.Infof("test:%d,%s,%t", 4, "zxc", true)
}
func TestLogFatalf(t *testing.T) {
	fl := FakeLogger{}
	fl.Fatalf("test:%d,%s,%t", 4, "zxc", true)
}
