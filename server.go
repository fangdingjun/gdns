/*
gdns is a dns proxy server write by go.

gdns much like dnsmasq or chinadns, but it can run on windows.

Features:

    support different domains use different upstream dns servers
    support contact to the upstream dns server by tcp or udp
    support blacklist list to block the fake ip

Usage:

generate a config file and edit it
    $ gdns -dumpflags > dns.ini


run it
    $ sudo gdns -config dns.ini


*/
package main

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/miekg/dns"
	"log"
	"strings"
)

var client_udp *dns.Client = &dns.Client{}

var client_tcp *dns.Client = &dns.Client{Net: "tcp"}

var Servers []*UpstreamServer = nil

var logger *LogOut = nil

var Blacklist_ips Kv = nil

var debug bool = false

var dns_cache *lru.Cache

var hostfile string = ""
var record_hosts Hosts = nil

func in_blacklist(m *dns.Msg) bool {
	if Blacklist_ips == nil {
		return false
	}

	if m == nil {
		return false
	}

	for _, rr := range m.Answer {
		/* A */
		if t, ok := rr.(*dns.A); ok {
			ip := t.A.String()
			if _, ok1 := Blacklist_ips[ip]; ok1 {
				logger.Debug("%s is in blacklist\n", ip)
				return true
			}
		}

		/* AAAA */
		if t, ok := rr.(*dns.AAAA); ok {
			ip := t.AAAA.String()
			if _, ok1 := Blacklist_ips[ip]; ok1 {
				logger.Debug("%s is in blacklist\n", ip)
				return true
			}
		}
	}

	return false
}

func handleRoot(w dns.ResponseWriter, r *dns.Msg) {
	var err error
	var res *dns.Msg
	domain := r.Question[0].Name

	var done int

	/*
	   reply from hosts
	*/
	if record_hosts != nil {
		rr := record_hosts.Get(domain, r.Question[0].Qtype)
		if rr != nil {
			msg := new(dns.Msg)
			msg.SetReply(r)
			msg.Answer = append(msg.Answer, rr)
			w.WriteMsg(msg)
			logger.Debug("%s query %s %s %s, reply from hosts\n",
				w.RemoteAddr(),
				domain,
				dns.ClassToString[r.Question[0].Qclass],
				dns.TypeToString[r.Question[0].Qtype],
			)
			return
		}
	}

	key := fmt.Sprintf("%s_%s", domain, dns.TypeToString[r.Question[0].Qtype])

	// reply from cache
	if a, ok := dns_cache.Get(key); ok {
		msg := new(dns.Msg)
		msg.SetReply(r)

		aa := strings.Split(a.(string), "|")
		for _, a1 := range aa {
			rr, _ := dns.NewRR(a1)
			if rr != nil {
				msg.Answer = append(msg.Answer, rr)
			}
		}

		w.WriteMsg(msg)
		logger.Debug("%s query %s %s %s, reply from cache\n",
			w.RemoteAddr(),
			domain,
			dns.ClassToString[r.Question[0].Qclass],
			dns.TypeToString[r.Question[0].Qtype],
		)
		return
	}

	// forward to upstream server
	for i := 0; i < 2; i++ {
		done = 0
		for _, sv := range Servers {
			if sv.match(domain) {

				res, err = sv.query(r)
				if err != nil {
					logger.Error("%s", err)
					continue
				}

				logger.Debug("%s query %s %s %s, forward to %s:%s, %s\n",
					w.RemoteAddr(),
					domain,
					dns.ClassToString[r.Question[0].Qclass],
					dns.TypeToString[r.Question[0].Qtype],
					sv.Proto, sv.Addr,
					dns.RcodeToString[res.Rcode],
				)

				if res.Rcode != dns.RcodeServerFailure && !in_blacklist(res) {
					// add to cache
					v := []string{}
					for _, as := range res.Answer {
						v = append(v, as.String())
					}
					dns_cache.Add(key, strings.Join(v, "|"))
					w.WriteMsg(res)
					done = 1
					break
				}
			}
		}

		// fallback to default upstream server
		if done != 1 {
			for _, dfsrv := range DefaultServer {
				res, err = dfsrv.query(r)
				if err != nil {
					logger.Error("%s", err)
					continue
				}

				logger.Debug("%s query %s %s %s, use default server %s:%s, %s\n",
					w.RemoteAddr(),
					domain,
					dns.ClassToString[r.Question[0].Qclass],
					dns.TypeToString[r.Question[0].Qtype],
					dfsrv.Proto, dfsrv.Addr,
					dns.RcodeToString[res.Rcode],
				)

				if res.Rcode != dns.RcodeServerFailure && !in_blacklist(res) {
					// add to cache
					v := []string{}
					for _, as := range res.Answer {
						v = append(v, as.String())
					}
					dns_cache.Add(key, strings.Join(v, "|"))
					w.WriteMsg(res)
					done = 1
					break
				}
			}
		}

		if done == 1 {
			break
		}
	}

	if done != 1 {
		dns.HandleFailed(w, r)
	}
}

func main() {
	parse_flags()

	var err error

	// create cache
	dns_cache, err = lru.New(1000)
	if err != nil {
		log.Fatal(err)
	}

	dns.HandleFunc(".", handleRoot)

	logger = NewLogger(logfile, debug)

	logger.Info("Listen on %s\n", bind_addr)

	go func() {
		/* listen tcp */
		err := dns.ListenAndServe(bind_addr, "tcp", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	/* listen udp */
	err = dns.ListenAndServe(bind_addr, "udp", nil)
	if err != nil {
		log.Fatal(err)
	}
}
