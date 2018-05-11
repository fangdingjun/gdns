package main

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
)

func TestCFDns(t *testing.T) {
	cf := &CloudflareHTTPDns{}
	m := new(dns.Msg)
	m.SetQuestion("www.google.com.", dns.TypeA)
	m1, _, err := cf.Exchange(m, "1.1.1.1")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(m1.String())
}
