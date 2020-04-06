/**
 * logger: manage all things log-related
 */

package logger

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/spf13/viper"
)

const (
	_ uint8 = iota
	// LogLevelNone will print no logs
	LogLevelNone
	// LogLevelPrint will print all logs up to PRINT level
	LogLevelPrint
	// LogLevelError will print all logs up to ERROR level
	LogLevelError
	// LogLevelWarn will print all logs up to WARN level
	LogLevelWarn
	// LogLevelInfo will print all logs up to INO level
	LogLevelInfo
	// LogLevelDebug will print all logs up to DEBUG level
	LogLevelDebug
	// LogLevelAll will print all logs
	LogLevelAll = ^uint8(0)
)

var (
	debuglog *log.Logger
	loglevel *uint8
)

func checkLogLevel() error {
	if loglevel == nil || *loglevel == uint8(0) {
		return fmt.Errorf("loglevel is not set")
	}
	return nil
}

// InitLogger initializes the logger
func InitLogger(initialLoglevel *uint8) error {
	debuglog = log.New(os.Stderr, "", log.LstdFlags|log.LUTC)
	l := uint8(viper.GetInt("loglevel"))

	if initialLoglevel != nil {
		l = *initialLoglevel
	}
	loglevel = &l

	return checkLogLevel()
}

// ResetLogLevel resets the global log level
func ResetLogLevel() error {
	l := uint8(viper.GetInt("loglevel"))
	loglevel = &l
	return checkLogLevel()
}

// Logger outputs a log to stdout
func logger(lvl uint8, fmt string, v ...interface{}) {
	if *loglevel >= lvl {
		debuglog.Printf(fmt, v...)
	}
	if lvl == LogLevelError {
		debuglog.Print(string(debug.Stack()))
	}
}

// Print prints a print level log
func Print(fmt string, v ...interface{}) {
	logger(LogLevelPrint, "[PRINT] "+fmt, v...)
}

// Error prints an error level log
func Error(fmt string, v ...interface{}) {
	logger(LogLevelError, "\033[31m[ERROR] \033[0m"+fmt, v...)
}

// Warn prints a warning level log
func Warn(fmt string, v ...interface{}) {
	logger(LogLevelWarn, "\033[33m[WARN ] \033[0m"+fmt, v...)
}

// Info prints an info level log
func Info(fmt string, v ...interface{}) {
	logger(LogLevelInfo, "\033[32m[INFO ] \033[0m"+fmt, v...)
}

// Debug prints a debug level log
func Debug(fmt string, v ...interface{}) {
	logger(LogLevelDebug, "[DEBUG] "+fmt, v...)
}

// Die outputs an error message before exiting
func Die(msg error) {
	debuglog.Fatalln("\033[31mfatal error:\033[0m", msg.Error())
}
