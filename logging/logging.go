package logging

import (
	"fmt"
	"log"
)

// Logger XXX
type Logger struct {
	tag string
}

type logLevel struct {
	name  string
	level uint
}

var trace = &logLevel{name: "TRACE", level: 1}
var debug = &logLevel{name: "DEBUG", level: 2}
var info = &logLevel{name: "INFO", level: 3}
var warning = &logLevel{name: "WARNING", level: 4}
var error = &logLevel{name: "ERROR", level: 5}
var critical = &logLevel{name: "CRITICAL", level: 6}

func stringTologLevel(name string) *logLevel {
	switch name {
	case "TRACE":
		return trace
	case "DEBUG":
		return debug
	case "INFO":
		return info
	case "WARNING":
		return warning
	case "CRITICAL":
		return critical
	}
	return &logLevel{name: "unknown", level: 0}
}

var logLevelConfigs = map[string]*logLevel{
	"root": info,
}

// GetLogger XXX
func GetLogger(tag string) *Logger {
	fmt.Printf("print something")
	return &Logger{tag: tag}
}

// ConfigureLoggers XXX
func ConfigureLoggers(rootlogLevel string) {
	logLevelConfigs["root"] = stringTologLevel(rootlogLevel)
}

func (logger *Logger) currentLogLevel() *logLevel {
	return logLevelConfigs["root"]
}

func (logger *Logger) message(logLevel *logLevel, message string) string {
	return logLevel.name + " " + logger.tag + " " + message
}

func (logger *Logger) log(logLevel *logLevel, message string, args ...interface{}) {
	if logLevel.level >= logger.currentLogLevel().level {
		log.Printf(logger.message(logLevel, message), args...)
	}
}

// Critical XXX
func (logger *Logger) Critical(m string, args ...interface{}) {
	logger.log(critical, m, args...)
}

// Error XXX
func (logger *Logger) Error(m string, args ...interface{}) {
	logger.log(error, m, args...)
}

// Warning XXX
func (logger *Logger) Warning(m string, args ...interface{}) {
	logger.log(warning, m, args...)
}

// Info XXX
func (logger *Logger) Info(m string, args ...interface{}) {
	logger.log(info, m, args...)
}

// Debug XXX
func (logger *Logger) Debug(m string, args ...interface{}) {
	logger.log(debug, m, args...)
}

// Trace XXX
func (logger *Logger) Trace(m string, args ...interface{}) {
	logger.log(trace, m, args...)
}
