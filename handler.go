package main

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"strings"
	"time"
)

type dnsClient interface {
	Exchange(m *dns.Msg, addr string) (*dns.Msg, time.Duration, error)
}

type dnsHandler struct {
	cfg         *conf
	tcpclient   dnsClient
	udpclient   dnsClient
	httpsclient dnsClient
}

func newDNSHandler(cfg *conf) *dnsHandler {
	return &dnsHandler{
		cfg:         cfg,
		tcpclient:   &dns.Client{Net: "tcp", Timeout: 2 * time.Second},
		udpclient:   &dns.Client{Net: "udp", Timeout: 2 * time.Second},
		httpsclient: &GoogleHTTPDns{},
	}

}

// ServerDNS implements the dns.Handler interface
func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	domain := r.Question[0].Name

	if ok := h.answerFromHosts(w, r); ok {
		return
	}

	srvs := h.getUpstreamServer(domain)
	if srvs == nil {
		srvs = h.cfg.DefaultUpstream
	}

	if msg, err := h.getAnswerFromUpstream(r, srvs); err == nil {
		w.WriteMsg(msg)
		return
	}

	dns.HandleFailed(w, r)

}

func (h *dnsHandler) getUpstreamServer(domain string) []addr {
	for _, srv := range h.cfg.ForwardRules {
		if ok := srv.domains.has(strings.Trim(domain, ".")); ok {
			return srv.Server
		}
	}
	return nil
}

func (h *dnsHandler) queryUpstream(r *dns.Msg, srv addr, ch chan *dns.Msg) {
	var m *dns.Msg
	var err error

	switch srv.Network {
	case "tcp":
		info("query %s IN %s, forward to %s:%d through tcp",
			r.Question[0].Name,
			dns.TypeToString[r.Question[0].Qtype],
			srv.Host,
			srv.Port)
		m, _, err = h.tcpclient.Exchange(r, fmt.Sprintf("%s:%d", srv.Host, srv.Port))
	case "udp":
		info("query %s IN %s, forward to %s:%d through udp",
			r.Question[0].Name,
			dns.TypeToString[r.Question[0].Qtype],
			srv.Host,
			srv.Port)
		m, _, err = h.udpclient.Exchange(r, fmt.Sprintf("%s:%d", srv.Host, srv.Port))
	case "https":
		info("query %s IN %s, forward to %s:%d through https",
			r.Question[0].Name,
			dns.TypeToString[r.Question[0].Qtype],
			srv.Host,
			srv.Port)
		m, _, err = h.httpsclient.Exchange(r, fmt.Sprintf("%s:%d", srv.Host, srv.Port))
	default:
		// ignore
	}

	if err == nil {
		select {
		case ch <- m:
		default:
		}
	} else {
		errorlog("%s", err)
	}
}

func (h *dnsHandler) getAnswerFromUpstream(r *dns.Msg, servers []addr) (*dns.Msg, error) {
	ch := make(chan *dns.Msg, 5)
	for _, srv := range servers {
		go func(a addr) {
			h.queryUpstream(r, a, ch)
		}(srv)
	}

	var savedErr *dns.Msg
	for {
		select {
		case m := <-ch:
			if m.Rcode == dns.RcodeSuccess && !h.inBlacklist(m) {
				return m, nil
			}
			savedErr = m
		case <-time.After(time.Duration(h.cfg.Timeout) * time.Second):
			if savedErr != nil {
				return savedErr, nil
			}
			info("query %s IN %s, timeout", r.Question[0].Name, dns.TypeToString[r.Question[0].Qtype])
			return nil, errors.New("timeout")
		}
	}
}

func (h *dnsHandler) inBlacklist(m *dns.Msg) bool {
	var ip string
	for _, rr := range m.Answer {
		if a, ok := rr.(*dns.A); ok {
			ip = a.String()
		} else if aaaa, ok := rr.(*dns.AAAA); ok {
			ip = aaaa.String()
		} else {
			ip = ""
		}
		if ip != "" && h.cfg.blacklist.exists(ip) {
			info("%s in blacklist", ip)
			return true
		}
	}
	return false
}

func (h *dnsHandler) answerFromHosts(w dns.ResponseWriter, r *dns.Msg) bool {
	domain := r.Question[0].Name
	t := r.Question[0].Qtype

	ip := h.cfg.hosts.get(strings.Trim(domain, "."), int(t))
	if ip != "" {
		rr, _ := dns.NewRR(fmt.Sprintf("%s 3600 IN %s %s", domain, dns.TypeToString[t], ip))
		if rr == nil {
			return false
		}
		msg := new(dns.Msg)
		msg.SetReply(r)
		msg.Answer = append(msg.Answer, rr)
		w.WriteMsg(msg)
		info("query %s IN %s, reply from hosts", domain, dns.TypeToString[t])
		return true
	}
	return false
}
