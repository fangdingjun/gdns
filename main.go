package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fangdingjun/go-log"
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
		log.Default.Out = &log.FixedSizeFileWriter{
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
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-ch:
		log.Errorf("received signal %s, exit...", s)
	}
	log.Println("exit.")
}
