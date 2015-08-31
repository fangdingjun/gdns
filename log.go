package main

import (
	"log"
	"os"
)

//var LogLevel int

type LogOut struct {
	//out *os.File
	debug   bool
	dbglog  *log.Logger
	errlog  *log.Logger
	infolog *log.Logger
}

func NewLogger(logfile string, debug bool) *LogOut {
	var out *os.File
	var err error
	if logfile != "" {
		out, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Println(err)
			out = os.Stdout
		}
	} else {
		out = os.Stdout
	}

	return &LogOut{
		debug,
		log.New(out, "DEBUG: ", log.LstdFlags),
		log.New(out, "ERROR: ", log.LstdFlags),
		log.New(out, "INFO: ", log.LstdFlags),
	}
}

func (l *LogOut) Debug(format string, args ...interface{}) {
	if l.debug {
		l.dbglog.Printf(format, args...)
	}
}

func (l *LogOut) Error(format string, args ...interface{}) {
	l.errlog.Printf(format, args...)
}

func (l *LogOut) Info(format string, args ...interface{}) {
	l.infolog.Printf(format, args...)
}
