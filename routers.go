package main

import (
	"github.com/miekg/dns"
	"log"
	"strings"
	"time"
)

type routers struct {
	c   *cfg
	tcp *dns.Client
	udp *dns.Client
}

func (r routers) checkBlacklist(m *dns.Msg) bool {
	if m.Rcode != dns.RcodeSuccess {
		// not success, not in blacklist
		return false
	}

	for _, rr := range m.Answer {
		var ip = ""
		if t, ok := rr.(*dns.A); ok {
			ip = t.A.String()
		} else if t, ok := rr.(*dns.AAAA); ok {
			ip = t.AAAA.String()
		}

		if ip != "" && r.c.blacklistIps.has(ip) {
			log.Printf("%s is in blacklist.\n", ip)
			return true
		}

	}
	return false
}

func (r routers) query(m *dns.Msg, servers []addr) (*dns.Msg, error) {
	var up *dns.Client
	var lastErr error
	for _, srv := range servers {
		switch srv.network {
		case "tcp":
			up = r.tcp
		case "udp":
			up = r.udp
		default:
			up = r.udp
		}

		log.Printf("query %s use %s:%s\n", m.Question[0].Name, srv.network, srv.addr)

		m, _, err := up.Exchange(m, srv.addr)
		if err == nil && !r.checkBlacklist(m) {
			return m, err
		}

		log.Println(err)
		lastErr = err
	}

	// return last error
	return nil, lastErr
}

// ServeDNS implements dns.Handler interface
func (r routers) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	domain := m.Question[0].Name
	d := strings.Trim(domain, ".")
	for _, rule := range r.c.Rules {
		if rule.domains.match(d) {
			m1, err := r.query(m, rule.servers)
			if err == nil {
				w.WriteMsg(m1)
				return
			} else {
				log.Println(err)
			}
		}
	}

	// no match or failed, fallback to default
	m1, err := r.query(m, r.c.servers)
	if err != nil {
		log.Println(err)
		dns.HandleFailed(w, m)
	} else {
		w.WriteMsg(m1)
	}
}

func initRouters(c *cfg) {
	router := &routers{
		c,
		&dns.Client{Net: "tcp", Timeout: time.Duration(c.Timeout) * time.Second},
		&dns.Client{Net: "udp", Timeout: time.Duration(c.Timeout) * time.Second},
	}
	dns.Handle(".", router)
}
