package slogger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDefaultLogLevel(t *testing.T) {
	if logLevel := GetLogLevel(); logLevel != infoLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := GetLogLevel(); logLevel == debugLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}

func TestSetLogLevelDebug(t *testing.T) {
	SetLogLevelDebug()
	if logLevel := GetLogLevel(); logLevel == infoLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := GetLogLevel(); logLevel != debugLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}

func TestSetLogLevelInfo(t *testing.T) {
	SetLogLevelInfo()
	if logLevel := GetLogLevel(); logLevel != infoLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := GetLogLevel(); logLevel == debugLogLevel {
		t.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}

func TestDebugMessage(t *testing.T) {
	SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	Debug("Test Message")
	w.Close()
	out, _ := ioutil.ReadAll(r)
	if ! strings.HasPrefix(string(out), "DEBUG") {
		t.Errorf("The Debug message didn't was not found on StdOut.")
	}
}

func TestDebugInfoMessage(t *testing.T) {
	SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	Debug("Test Message")
	w.Close()
	out, _ := ioutil.ReadAll(r)
	if ! strings.HasPrefix(string(out), "DEBUG") {
		t.Errorf("The Debug message didn't was not found on StdOut.")
	}
	r, w, _ = os.Pipe()
	os.Stdout = w
	Info("Test Message")
	w.Close()
	out, _ = ioutil.ReadAll(r)
	if ! strings.HasPrefix(string(out), "INFO") {
		t.Errorf("The Info message didn't was not found on StdOut.")
	}
}

func TestDebugMethod(t *testing.T) {
	SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	DebugMethod()
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
