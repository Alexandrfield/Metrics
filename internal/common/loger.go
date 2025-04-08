package common

type Loger interface {
	Warnf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}
type FakeLogger struct{}

func (fakeLogger *FakeLogger) Warnf(template string, args ...interface{})  {}
func (fakeLogger *FakeLogger) Infof(template string, args ...interface{})  {}
func (fakeLogger *FakeLogger) Fatalf(template string, args ...interface{}) {}
func (fakeLogger *FakeLogger) Errorf(template string, args ...interface{}) {}
func (fakeLogger *FakeLogger) Debugf(template string, args ...interface{}) {}
