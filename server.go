package main

import (
	"flag"
	"github.com/miekg/dns"
	"log"
	"os"
)

func initListeners(c *cfg) {
	for _, a := range c.listen {
		log.Printf("Listen on %s %s...\n", a.network, a.addr)
		s := dns.Server{Addr: a.addr, Net: a.network}
		go s.ListenAndServe()
	}
}

func main() {
	var configFile string

	flag.StringVar(&configFile, "c", "", "config file")
	flag.Parse()

	config, err := parseCfg(configFile)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	initRouters(config)
	initListeners(config)

	select {}
}
