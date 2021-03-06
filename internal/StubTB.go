package internal

import (
	"runtime"
	"testing"
)

type StubTB struct {
	// TB is only present here to implement testing.TB interface's
	// unexported functions by embedding the interface itself.
	testing.TB

	IsFailed  bool
	IsSkipped bool

	StubName    string
	StubTempDir string
	StubCleanup func(func())

	cleanups []func()
}

func (m *StubTB) Finish() {
	for _, fn := range m.cleanups {
		defer fn()
	}
}

func (m *StubTB) Cleanup(f func()) {
	fn := func() { InGoroutine(f) }
	if m.StubCleanup != nil {
		m.StubCleanup(fn)
	} else {
		m.cleanups = append(m.cleanups, fn)
	}
}

func (m *StubTB) Error(args ...interface{}) {
	m.Fail()
}

func (m *StubTB) Errorf(format string, args ...interface{}) {
	m.Fail()
}

func (m *StubTB) Fail() {
	m.IsFailed = true
}

func (m *StubTB) FailNow() {
	m.Fail()
	runtime.Goexit()
}

func (m *StubTB) Failed() bool {
	return m.IsFailed
}

func (m *StubTB) Fatal(args ...interface{}) {
	m.FailNow()
}

func (m *StubTB) Fatalf(format string, args ...interface{}) {
	m.FailNow()
}

func (m *StubTB) Helper() {}

func (m *StubTB) Log(args ...interface{}) {}

func (m *StubTB) Logf(format string, args ...interface{}) {}

func (m *StubTB) Name() string {
	return m.StubName
}

func (m *StubTB) Skip(args ...interface{}) {
	m.SkipNow()
}

func (m *StubTB) SkipNow() {
	m.IsSkipped = true
	runtime.Goexit()
}

func (m *StubTB) Skipf(format string, args ...interface{}) {
	m.SkipNow()
}

func (m *StubTB) Skipped() bool {
	return m.IsSkipped
}

func (m *StubTB) TempDir() string {
	return m.StubTempDir
}
