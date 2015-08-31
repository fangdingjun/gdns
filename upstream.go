package main

import (
	"github.com/miekg/dns"
	"strings"
)

type Kv map[string]int

type UpstreamServer struct {
	domains Kv
	Proto   string
	Addr    string
	client  *dns.Client
}

func (srv *UpstreamServer) match(d string) bool {
	if srv.domains == nil {
		return false
	}

	s := strings.Split(strings.Trim(d, "."), ".")

	for i := 0; i < len(s)-1; i++ {
		s1 := strings.Join(s[i:], ".")
		if _, ok := srv.domains[s1]; ok {
			return true
		}
	}

	return false
}

func (srv *UpstreamServer) query(req *dns.Msg) (*dns.Msg, error) {
	res, _, err := srv.client.Exchange(req, srv.Addr)
	return res, err
}
