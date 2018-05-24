package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fangdingjun/go-log"
	"github.com/fangdingjun/go-log/formatters"
	"github.com/fangdingjun/go-log/writers"
	"github.com/miekg/dns"
)

func main() {
	var configfile string
	var logfile string
	var loglevel string
	var logFileSize int64
	var logFileCount int

	flag.StringVar(&configfile, "c", "config.yaml", "config file")
	flag.StringVar(&logfile, "log_file", "", "log file, default stdout")
	flag.IntVar(&logFileCount, "log_count", 10, "max count of log to keep")
	flag.Int64Var(&logFileSize, "log_size", 10, "max log file size MB")
	flag.StringVar(&loglevel, "log_level", "INFO", "log level, values:\nOFF, FATAL, PANIC, ERROR, WARN, INFO, DEBUG")
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

	log.Default.Formatter = &formatters.TextFormatter{TimeFormat: "2006-01-02 15:04:05.000"}

	config, err := loadConfig(configfile)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("config: %+v", config)

	h := newDNSHandler(config)

	for _, l := range config.Listen {
		go func(l addr) {
			log.Infof("listen on %s %s:%d", l.Network, l.Host, l.Port)
			if err := dns.ListenAndServe(
				fmt.Sprintf("%s:%d", l.Host, l.Port), l.Network, h); err != nil {
				log.Fatal(err)
			}
		}(l)
	}

	select {}
}
