package utils

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type aggregatedLogger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

type ILogger interface {
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	FatalError(v ...interface{})
}

// ANSI escape codes for color
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

func NewLogger() ILogger {
	flags := log.LstdFlags

	return &aggregatedLogger{
		infoLogger:  log.New(os.Stdout, colorCyan+"INFO: "+colorReset, flags),
		warnLogger:  log.New(os.Stdout, colorYellow+"WARN: "+colorReset, flags),
		errorLogger: log.New(os.Stdout, colorRed+"ERROR: "+colorReset, flags),
	}
}

func (l *aggregatedLogger) logWithCaller(logger *log.Logger, v ...interface{}) {
	_, file, line, _ := runtime.Caller(2) // Adjust the caller depth based on your package structure
	_, fileName := filepath.Split(file)
	if fileName[0:6] == "logmid" {
		logger.Printf("%v", v)
	} else {
		logger.Printf("%s:%d %v", fileName, line, v)
	}
}

func (l *aggregatedLogger) Info(v ...interface{}) {
	l.logWithCaller(l.infoLogger, v...)
}

func (l *aggregatedLogger) Warn(v ...interface{}) {
	l.logWithCaller(l.warnLogger, v...)
}

func (l *aggregatedLogger) Error(v ...interface{}) {
	l.logWithCaller(l.errorLogger, v...)
}

func (l *aggregatedLogger) FatalError(v ...interface{}) {
	l.logWithCaller(l.errorLogger, v...)
	os.Exit(1)
}
