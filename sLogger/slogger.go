/*
sLogger is a wrapper for the log package for Sote GO software developers.

The GO log package should not be used directly.  This package sets the format of
the logging message so they are uniform.  Also, Sote only has two logging level,
DEBUG and INFO.  If the system has failed, then panic should be called after all
available information has been output to the log.

The log message format is:
    {year}/{month}}/{day} {hour}:{mins}:{secs}.{microsecs} {[logPrefix.]MessageType}:{Test Message}
    example:
    2020/06/16 22:26:42.165609 SLOGGER_TEST.DEBUG:gitlab.com/soteapps/packages/v2020/sLogger.TestSetLogMessagePrefix
It is recommended that the application set the log prefix (SetLogMessagePrefix) so log messages can be easily grouped.  If not, "missing" will be used.

*/
package sLogger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	DebugLogLevel    = "DEBUG"   // Support log levels
	InfoLogLevel     = "INFO"    // Support log levels
	logPrefixMissing = "missing" // Support log levels
)

var (
	logMessage *log.Logger
	logLevel   string = InfoLogLevel
	logPrefix  string = logPrefixMissing
)

// This is used to set the logging message format
func initLogger(infoHandle io.Writer, msgType string) {
	logMessage = log.New(infoHandle, fmt.Sprintf("%v.%v:", logPrefix, msgType), log.Lmsgprefix|log.LstdFlags|log.Lmicroseconds|log.LUTC)
}

// This will publish a log message at the INFO level
func Info(tMessage string) {
	initLogger(os.Stdout, InfoLogLevel)
	logMessage.Println(tMessage)
}

// This will publish a log message at the DEBUG level
func Debug(tMessage string) {
	if logLevel == DebugLogLevel {
		initLogger(os.Stdout, DebugLogLevel)
		logMessage.Println(tMessage)
	}
}

// This will publish a log message at the DEBUG level for the function that is being executed.
func DebugMethod(depthList ...int) {
	if logLevel == DebugLogLevel {
		initLogger(os.Stdout, DebugLogLevel)
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
	logLevel = DebugLogLevel
}

// This will set the log level to INFO
func SetLogLevelInfo() {
	logLevel = InfoLogLevel
}

// This will return the logging level (It doesn't follow return referring to the func declaration)
func GetLogLevel() string {
	return logLevel
}

// Allows a message prefix to be applied
// The prefix is forced to upper case
func SetLogMessagePrefix(prefix string) {
	logPrefix = strings.ToUpper(prefix)
}
