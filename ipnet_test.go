package main

import (
	//"fmt"
	"github.com/miekg/dns"
	"testing"
	"time"
	//"os"
)

func TestParseNet(t *testing.T) {
	nets := parse_net("region_cn.txt")
	t.Logf("get %d networks\n", len(nets))
	t.Logf("1st %s\n", nets[0].String())
}

func TestQuery(t *testing.T) {
	ip_region = parse_net("region_cn.txt")
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
		t.Log(err)
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
		t1 := time.Now()
		m1 := query(m)
		t2 := time.Now()
		if m1 == nil {
			t.Errorf("query %s failed", dn)
		} else {
			t.Logf("query time: %s\n", t2.Sub(t1))
			t.Logf("result of %s\n", dn)
			for _, a1 := range m1.Answer {
				t.Logf("%s\n", a1)
			}
			//.Printf("\n")
		}
	}
}
