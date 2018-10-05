package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fangdingjun/go-log"
	"github.com/fangdingjun/go-log/formatters"
	"github.com/fangdingjun/go-log/writers"
)

func main() {
	var configfile string
	var logFileCount int
	var logFileSize int64
	var loglevel string
	var logfile string

	flag.StringVar(&logfile, "log_file", "", "log file, default stdout")
	flag.IntVar(&logFileCount, "log_count", 10, "max count of log to keep")
	flag.Int64Var(&logFileSize, "log_size", 10, "max log file size MB")
	flag.StringVar(&loglevel, "log_level", "INFO",
		"log level, values:\nOFF, FATAL, PANIC, ERROR, WARN, INFO, DEBUG")
	flag.StringVar(&configfile, "c", "gdns.yaml", "config file")
	flag.Parse()

	if logfile != "" {
		log.Default.Out = &writers.FixedSizeFileWriter{
			MaxCount: logFileCount,
			Name:     logfile,
			MaxSize:  logFileSize * 1024 * 1024,
		}
	}

	if loglevel != "" {
		lv, err := log.ParseLevel(loglevel)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		log.Default.Level = lv
	}

	log.Default.Formatter = &formatters.TextFormatter{
		TimeFormat: "2006-01-02 15:04:05.000"}

	cfg, err := loadConfig(configfile)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.UpstreamTimeout == 0 {
		cfg.UpstreamTimeout = 5
	}
	initDNSClient(cfg)
	log.Debugf("%+v", cfg)
	makeServers(cfg)
	select {}
}
