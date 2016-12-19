package main

import (
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"log"
)

func main() {
	var configfile string

	flag.StringVar(&configfile, "c", "config.yaml", "config file")
	flag.Parse()

	config, err := loadConfig(configfile)
	if err != nil {
		log.Fatal(err)
	}

	h := newDNSHandler(config)

	for _, l := range config.Listen {
		go func(l addr) {
			if err := dns.ListenAndServe(
				fmt.Sprintf("%s:%d", l.Host, l.Port), l.Network, h); err != nil {
				log.Fatal(err)
			}
		}(l)
	}

	select {}
}
