package sLogger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDefaultLogLevel(tPtr *testing.T) {
	if logLevel := GetLogLevel(); logLevel != InfoLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := GetLogLevel(); logLevel == DebugLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}
func TestSetLogLevelDebug(tPtr *testing.T) {
	SetLogLevelDebug()
	if logLevel := GetLogLevel(); logLevel == InfoLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := GetLogLevel(); logLevel != DebugLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}
func TestSetLogLevelInfo(tPtr *testing.T) {
	SetLogLevelInfo()
	if logLevel := GetLogLevel(); logLevel != InfoLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := GetLogLevel(); logLevel == DebugLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}
func TestDebugMessage(tPtr *testing.T) {
	SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	Debug("Test Message")
	w.Close()
	out, _ := ioutil.ReadAll(r)
	if ! strings.Contains(string(out), DebugLogLevel) {
		tPtr.Errorf("The Debug message didn't was not found on StdOut.")
	}
}
func TestDebugInfoMessage(tPtr *testing.T) {
	SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	Debug("Test Message")
	w.Close()
	out, _ := ioutil.ReadAll(r)
	if ! strings.Contains(string(out), DebugLogLevel) {
		tPtr.Errorf("The Debug message didn't was not found on StdOut.")
	}
	r, w, _ = os.Pipe()
	os.Stdout = w
	Info("Test Message")
	w.Close()
	out, _ = ioutil.ReadAll(r)
	if ! strings.Contains(string(out), InfoLogLevel) {
		tPtr.Errorf("The Info message didn't was not found on StdOut.")
	}
}
func TestDebugMethod(tPtr *testing.T) {
	SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	DebugMethod()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	output := strings.TrimSuffix(string(out), "\n")
	if ! strings.Contains(output, DebugLogLevel) {
		tPtr.Errorf("The Debug message didn't was not found on StdOut.")
	}
	if ! strings.Contains(output, "TestDebugMethod") {
		tPtr.Errorf("The Debug message didn't have the func name for this test in Stdout.")
	}
}
func TestSetLogMessagePrefix(tPtr *testing.T) {
	SetLogLevelDebug()
	SetLogMessagePrefix("sLogger_test")
	r, w, _ := os.Pipe()
	os.Stdout = w
	DebugMethod()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	output := strings.TrimSuffix(string(out), "\n")
	if ! strings.Contains(output, DebugLogLevel) && ! strings.Contains(output, "SLOGGER_TEST") {
		tPtr.Errorf("The Debug message didn't was not found on StdOut.")
	}
	if ! strings.Contains(output, "TestDebugMethod") && ! strings.Contains(output, "SLOGGER_TEST") {
		tPtr.Errorf("The Debug message didn't have the func name for this test in Stdout.")
	}
}
