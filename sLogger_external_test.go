package packages

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"gitlab.com/soteapps/packages/v2021/sLogger"
)

func TestDefaultLogLevel(tPtr *testing.T) {
	if logLevel := sLogger.GetLogLevel(); logLevel != sLogger.InfoLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := sLogger.GetLogLevel(); logLevel == sLogger.DebugLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
}
func TestSetLogLevelDebug(tPtr *testing.T) {
	sLogger.SetLogLevelDebug()
	if logLevel := sLogger.GetLogLevel(); logLevel == sLogger.InfoLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := sLogger.GetLogLevel(); logLevel != sLogger.DebugLogLevel {
		tPtr.Errorf("The default log level should be Debug.  It is " + logLevel)
	}
}
func TestSetLogLevelInfo(tPtr *testing.T) {
	sLogger.SetLogLevelInfo()
	if logLevel := sLogger.GetLogLevel(); logLevel != sLogger.InfoLogLevel {
		tPtr.Errorf("The default log level should be INFO.  It is " + logLevel)
	}
	if logLevel := sLogger.GetLogLevel(); logLevel == sLogger.DebugLogLevel {
		tPtr.Errorf("The default log level should be Debug.  It is " + logLevel)
	}
}
func TestDebugMessage(tPtr *testing.T) {
	sLogger.SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	sLogger.Debug("Test Message")
	w.Close()
	out, _ := ioutil.ReadAll(r)
	if ! strings.Contains(string(out), sLogger.DebugLogLevel) {
		tPtr.Errorf("The Debug message was not found on StdOut.")
	}
}
func TestDebugInfoMessage(tPtr *testing.T) {
	sLogger.SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	sLogger.Debug("Test Message")
	w.Close()
	out, _ := ioutil.ReadAll(r)
	if ! strings.Contains(string(out), sLogger.DebugLogLevel) {
		tPtr.Errorf("The Debug message was not found on StdOut.")
	}
	r, w, _ = os.Pipe()
	os.Stdout = w
	sLogger.Info("Test Message")
	w.Close()
	out, _ = ioutil.ReadAll(r)
	if ! strings.Contains(string(out), sLogger.InfoLogLevel) {
		tPtr.Errorf("The Info message was not found on StdOut.")
	}
}
func TestDebugMethod(tPtr *testing.T) {
	sLogger.SetLogLevelDebug()
	r, w, _ := os.Pipe()
	os.Stdout = w
	sLogger.DebugMethod()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	output := strings.TrimSuffix(string(out), "\n")
	if ! strings.Contains(output, sLogger.DebugLogLevel) {
		tPtr.Errorf("The Debug message was not found on StdOut.")
	}
	if ! strings.Contains(output, "TestDebugMethod") {
		tPtr.Errorf("The Debug message have the func name for this test in Stdout.")
	}
}
func TestSetLogMessagePrefix(tPtr *testing.T) {
	sLogger.SetLogLevelDebug()
	sLogger.SetLogMessagePrefix("sLogger_test")
	r, w, _ := os.Pipe()
	os.Stdout = w
	sLogger.DebugMethod()
	w.Close()
	out, _ := ioutil.ReadAll(r)
	output := strings.TrimSuffix(string(out), "\n")
	if ! strings.Contains(output, sLogger.DebugLogLevel) && ! strings.Contains(output, "SLOGGER_TEST") {
		tPtr.Errorf("The Debug message didn't was not found on StdOut.")
	}
	if ! strings.Contains(output, "TestDebugMethod") && ! strings.Contains(output, "SLOGGER_TEST") {
		tPtr.Errorf("The Debug message didn't have the func name for this test in Stdout.")
	}
}
