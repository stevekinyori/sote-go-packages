package packages

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gitlab.com/getsote/packages/slogger"
)

func TestDefaultLogLevel(t *testing.T) {
	if logLevel := slogger.GetLogLevel(); logLevel != slogger.InfoLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := slogger.GetLogLevel(); logLevel == slogger.DebugLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}

func TestSetLogLevelDebug(t *testing.T) {
	slogger.SetLogLevelDebug()
	if logLevel := slogger.GetLogLevel(); logLevel == slogger.InfoLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := slogger.GetLogLevel(); logLevel != slogger.DebugLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}

func TestSetLogLevelInfo(t *testing.T) {
	slogger.SetLogLevelInfo()
	if logLevel := slogger.GetLogLevel(); logLevel != slogger.InfoLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := slogger.GetLogLevel(); logLevel == slogger.DebugLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}

func TestDebugMessage(t *testing.T) {
	slogger.SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	slogger.Debug("Test Message")
	w.Close()
	out, _ := ioutil.ReadAll(r)
	if ! strings.HasPrefix(string(out), "DEBUG") {
		t.Errorf("The Debug message didn't was not found on StdOut.")
	}
}

func TestDebugInfoMessage(t *testing.T) {
	slogger.SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	slogger.Debug("Test Message")
	w.Close()
	out, _ := ioutil.ReadAll(r)
	if ! strings.HasPrefix(string(out), "DEBUG") {
		t.Errorf("The Debug message didn't was not found on StdOut.")
	}
	r, w, _ = os.Pipe()
	os.Stdout = w
	slogger.Info("Test Message")
	w.Close()
	out, _ = ioutil.ReadAll(r)
	if ! strings.HasPrefix(string(out), "INFO") {
		t.Errorf("The Info message didn't was not found on StdOut.")
	}
}

func TestDebugMethod(t *testing.T) {
	slogger.SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	slogger.DebugMethod()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	output := strings.TrimSuffix(string(out), "\n")
	if ! strings.HasPrefix(output, "DEBUG") {
		t.Errorf("The Debug message didn't was not found on StdOut.")
	}
	if ! strings.HasSuffix(output, "TestDebugMethod") {
		t.Errorf("The Debug message didn't have the func name for this test in Stdout.")
	}
}
