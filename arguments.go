package main

import (
	"flag"
	"github.com/miekg/dns"
	"github.com/vharitonsky/iniflags"
	"log"
)

var bind_addr string

var default_server string

var srv ArgSrvs

var logfile string

type ArgSrvs []string

var DefaultServer *UpstreamServer

var blacklist_file string

func (s *ArgSrvs) String() string {
	//Sprintf("%s", s)
	return "filter1.txt,udp:8.8.8.8:53"
}

func (s *ArgSrvs) Set(s1 string) error {
	*s = append(*s, s1)
	return nil
}

func parse_flags() {
	iniflags.Parse()

	var err error
	for _, s := range srv {
		sv, err := parse_server(s)
		if err != nil {
			log.Print(err)
		} else {
			Servers = append(Servers, sv)
		}
	}

	proto, addr, err := parse_addr(default_server)
	if err != nil {
		log.Fatal(err)
	}

	var c *dns.Client
	if proto == "udp" {
		c = client_udp
	} else {
		c = client_tcp
	}

	DefaultServer = &UpstreamServer{
		Addr:   addr,
		Proto:  proto,
		client: c,
	}

	a, err := load_domain(blacklist_file)
	if err != nil {
		log.Println(err)
	} else {
		Blacklist_ips = a
	}
}

func init() {

	flag.Var(&srv, "server", `special the filter and the upstream server to use when match
      format:
          FILTER_FILE_NAME,PROTOCOL:SERVER_NAME:PORT
      example:
          filter1.json,udp:8.8.8.8:53
              means the domains in the filter1.json will use the google dns  server by udp
      you can specail multiple filter and upstream server        
        `)

	flag.StringVar(&bind_addr, "bind", ":53", "the address bind to")
	flag.StringVar(&default_server, "upstream", "udp:114.114.114.114:53", "the default upstream server to use")
	flag.StringVar(&logfile, "logfile", "error.log", "the logfile, default stdout")
	flag.StringVar(&blacklist_file, "blacklist", "", "the blacklist file")
	flag.BoolVar(&debug, "debug", false, "output debug log, default false")
}
