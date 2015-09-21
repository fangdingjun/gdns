package main

import (
	"github.com/miekg/dns"
	"testing"
	//"fmt"
)

func TestReadHosts(t *testing.T) {
	a, err := ReadHosts("testdata/hosts.txt")
	if err != nil {
		t.Error(err)
	}

	for k, v := range a {
		for _, v1 := range v {
			t.Logf("%s: %s\n", k, v1.rr.String())
		}
	}
	r1 := a.Get("localhost", dns.TypeA)
	if dnsa, ok := r1.(*dns.A); ok {
		if dnsa.A.String() != "127.0.0.1" {
			t.Errorf("get failed a\n")
		}
	} else {
		t.Errorf("type not a\n")
	}
	r2 := a.Get("localhost", dns.TypeAAAA)
	if dnsaa, ok := r2.(*dns.AAAA); ok {
		if dnsaa.AAAA.String() != "::1" {
			t.Errorf("get failed aaaa\n")
		}
	} else {
		t.Errorf("type not aaaa\n")
	}
}

func TestGetHost(t *testing.T) {
	t.Logf("host: %s\n", GetHost())
}
