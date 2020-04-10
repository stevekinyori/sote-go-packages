// This is a wrapper for the log package for Sote GO software developers.
//
// The GO log package should not be used directly.  This package sets the format of
// the logging message so they are uniform.  Also, Sote only has two logging level,
// DEBUG and INFO.  If the system has failed, then panic should be called after all
// available information has been output to the log.
package slogger

import (
	"io"
	"log"
	"os"
	"runtime"
)

const debugLogLevel = "DEBUG" // Support log levels
const infoLogLevel = "INFO"   // Support log levels

var (
	logMessage *log.Logger
	logLevel   string = infoLogLevel
)

// This is used to set the logging message format
func initLogger(infoHandle io.Writer, msgType string) {
	logMessage = log.New(infoHandle, msgType, log.LstdFlags|log.LUTC)
}

// This will publish a log message at the INFO level
func Info(tMessage string) {
	initLogger(os.Stdout, "INFO:")
	logMessage.Println(tMessage)
}

// This will publish a log message at the DEBUG level
func Debug(tMessage string) {
	if logLevel == debugLogLevel {
		initLogger(os.Stdout, "DEBUG:")
		logMessage.Println(tMessage)
	}
}

// This will publish a log message at the DEBUG level for the function that is being executed.
func DebugMethod(depthList ...int) {
	if logLevel == debugLogLevel {
		initLogger(os.Stdout, "DEBUG:")
		var depth int
		if depthList == nil {
			depth = 1
		} else {
			depth = depthList[0]
		}
		function, _, _, _ := runtime.Caller(depth)
		functionName := runtime.FuncForPC(function).Name()
		logMessage.Println(functionName)
	}
}

// This will set the log level to DEBUG
func SetLogLevelDebug() {
	logLevel = "DEBUG"
}

// This will set the log level to INFO
func SetLogLevelInfo() {
	logLevel = "INFO"
}

// This will return the logging level
func GetLogLevel() string {
	return logLevel
}
