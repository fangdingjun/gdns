package main

import (
	"errors"
	"github.com/miekg/dns"
	"log"
	"strings"
	"time"
)

type routers struct {
	c     *cfg
	tcp   *dns.Client
	udp   *dns.Client
	cache *cache
}

func (r *routers) checkBlacklist(m *dns.Msg) bool {
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

type dnsClient interface {
	Exchange(*dns.Msg, string) (*dns.Msg, time.Duration, error)
}

func (r *routers) query(m *dns.Msg, servers []addr) (*dns.Msg, error) {
	var up dnsClient
	var lastErr error

	// query cache
	m2 := r.cache.get(m)
	if m2 != nil {
		log.Printf("query %s, reply from cache\n", m.Question[0].Name)
		m2.Id = m.Id
		return m2, nil
	}

	for _, srv := range servers {
		switch srv.network {
		case "tcp":
			up = r.tcp
		case "udp":
			up = r.udp
		case "https":
			up = DefaultHTTPDnsClient
		default:
			up = r.udp
		}

		log.Printf("query %s use %s:%s\n", m.Question[0].Name, srv.network, srv.addr)

		m1, _, err := up.Exchange(m, srv.addr)
		if err == nil && !r.checkBlacklist(m) {
			if m1.Rcode == dns.RcodeSuccess {
				// store to cache
				r.cache.set(m1)
			}
			return m1, err
		}

		log.Println(err)
		lastErr = err
	}

	if lastErr == nil {
		// this happens when ip in blacklist
		lastErr = errors.New("timeout")
	}

	// return last error
	return nil, lastErr
}

// ServeDNS implements dns.Handler interface
func (r *routers) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	domain := m.Question[0].Name
	d := strings.Trim(domain, ".")
	for _, rule := range r.c.Rules {
		if rule.domains.match(d) {
			m1, err := r.query(m, rule.servers)
			if err == nil {
				w.WriteMsg(m1)
				return
			}

			log.Println(err)

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
		newCache(1000, int64(c.TTL)), // cache 5 hours
	}
	dns.Handle(".", router)
}
