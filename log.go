package main

import (
	"log"
)

const (
	_ = iota
	FATAL
	ERROR
	WARN
	NOTICE
	INFO
	DEBUG
)

var logLevel = WARN

func logMsg(l int, fmt string, args ...interface{}) {
	if l <= logLevel {
		log.Printf(fmt, args...)
	}
}

func info(fmt string, args ...interface{}) {
	logMsg(INFO, fmt, args...)
}

func errorlog(fmt string, args ...interface{}) {
	logMsg(ERROR, fmt, args...)
}
