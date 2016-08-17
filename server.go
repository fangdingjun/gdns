package main

import (
	"flag"
	"github.com/fangdingjun/gpp/util"
	"github.com/miekg/dns"
	"log"
	"os"
	"time"
)

func initListeners(c *cfg) {
	for _, a := range c.listen {
		log.Printf("Listen on %s %s...\n", a.network, a.addr)
		s := &dns.Server{Addr: a.addr, Net: a.network}
		go func(s *dns.Server) {
			err := s.ListenAndServe()
			if err != nil {
				log.Println(err)
				os.Exit(-1)
			}
		}(s)
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

	// make a delay to make sure net bind completed before drop privilege
	time.Sleep(time.Second)

	err = util.DropPrivilege(config.User, config.Group)
	if err != nil {
		log.Println(err)
	}

	select {}
}
