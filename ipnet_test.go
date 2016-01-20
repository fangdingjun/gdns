package main

import (
	"fmt"
	"github.com/miekg/dns"
	"testing"
	//"os"
)

func TestParseNet(t *testing.T) {
	nets := parse_net("net_china.txt")
	fmt.Printf("get %d networks\n", len(nets))
	fmt.Printf("1st %s\n", nets[0].String())
}

func TestQuery(t *testing.T) {
	ip_region = parse_net("net_china.txt")
	var c *dns.Client
	for _, srv := range []string{
		"tcp:114.114.114.114:53",
		"udp:8.8.8.8:53",
		"udp:192.168.41.1:53",
		"udp:4.2.2.2:53",
	} {
		proto, addr, err := parse_addr(srv)
		if err != nil {
			t.Error(err)
		}
		if proto == "tcp" {
			c = client_tcp
		} else {
			c = client_udp
		}
		upsrv := &UpstreamServer{
			Addr:   addr,
			Proto:  proto,
			client: c,
		}
		DefaultServer = append(DefaultServer, upsrv)
	}
	blacklist_file = "blacklist.txt"
	a, err := load_domain(blacklist_file)
	if err != nil {
		fmt.Println(err)
	} else {
		Blacklist_ips = a
	}
	logger = NewLogger("", true)
	for _, dn := range []string{
		"www.google.com",
		"www.sina.com.cn",
		"www.taobao.com",
		"www.ifeng.com",
		"twitter.com",
		"www.facebook.com",
		"plus.google.com",
		"drive.google.com",
		"dongtaiwang.com",
		"www.ratafee.nl",
		"cc.ratafee.nl",
		"noddcade.xx.ffs.aafde",
		"ndfddcade.xx.ffs.aafde",
		"sddf32dsf.comd.ffdasdf.fdsd3eaaaaa",
		"www.google.com.hk",
	} {
		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(dn), dns.TypeA)
		m1 := query(m)
		if m1 == nil {
			t.Error("query failed")
		}
		fmt.Printf("%s\n", m1)
	}
}
