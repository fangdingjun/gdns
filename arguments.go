package main

import (
	"flag"
	"github.com/miekg/dns"
	"github.com/vharitonsky/iniflags"
	"log"
)

var bind_addr string

var default_server ArgSrvs

var srv ArgSrvs

var logfile string

type ArgSrvs []string

var DefaultServer []*UpstreamServer

var blacklist_file string

var enable_cache = false

var region_file = ""

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

	for _, dsvr := range default_server {
		proto, addr, err := parse_addr(dsvr)
		if err != nil {
			log.Fatal(err)
		}

		var c *dns.Client
		if proto == "udp" {
			c = client_udp
		} else {
			c = client_tcp
		}

		upsrv := &UpstreamServer{
			Addr:   addr,
			Proto:  proto,
			client: c,
		}
		DefaultServer = append(DefaultServer, upsrv)
	}

	if len(DefaultServer) == 0 {
		log.Fatal("please special a -upstream")
	}

	a, err := load_domain(blacklist_file)
	if err != nil {
		log.Println(err)
	} else {
		Blacklist_ips = a
	}

	if hostfile == "" {
		hostfile = GetHost()
	}

	if hostfile != "" {
		record_hosts, err = ReadHosts(hostfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if region_file != "" {
		ip_region = parse_net(region_file)
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
	flag.Var(&default_server, "upstream", "special the upstream server to use")
	flag.StringVar(&logfile, "logfile", "", "the logfile, default stdout")
	flag.StringVar(&blacklist_file, "blacklist", "", "the blacklist file")
	flag.BoolVar(&debug, "debug", false, "output debug log, default false")
	flag.StringVar(&hostfile, "hosts", "", "load special ip from hosts or /etc/hosts")
	flag.BoolVar(&enable_cache, "enable_cache", false, "enable cache or not")
	flag.StringVar(&region_file, "region_file", "", "local country region ip range file")
}
