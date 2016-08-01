package main

import (
	"github.com/miekg/dns"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := newCache(5, 2)

	tests := map[string]uint16{
		"www.google.com":    dns.TypeA,
		"www.google.com.hk": dns.TypeA,
		"www.google.com.sg": dns.TypeA,
		"www.google.com.it": dns.TypeA,
		"www.google.com.de": dns.TypeA,
		"www.google.com.cn": dns.TypeA,
	}

	var datas []*dns.Msg

	for k, v := range tests {
		m1 := new(dns.Msg)
		m1.SetQuestion(k, v)
		datas = append(datas, m1)
	}

	for i := 0; i < 3; i++ {
		c.set(datas[i])
	}

	for i := 0; i < 3; i++ {
		m2 := c.get(datas[i])
		if m2 == nil {
			t.Errorf("store cache failed")
		}
		if m2.Question[0].Name != datas[i].Question[0].Name {
			t.Errorf("cache error")
		}
	}

	time.Sleep(3 * time.Second)
	for i := 0; i < 3; i++ {
		m2 := c.get(datas[i])
		if m2 != nil {
			t.Errorf("cache not expired")
		}
	}

	for i := 3; i < 6; i++ {
		c.set(datas[i])
	}

	if len(c.m) > len(datas) {
		t.Errorf("old cache not purged")
	}
}
