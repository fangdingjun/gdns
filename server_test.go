package main

import (
	"fmt"
	"github.com/miekg/dns"
	"testing"
)

func TestInblacklist(t *testing.T) {
	logger = NewLogger("", false)
	Blacklist_ips = Kv{"1.2.3.4": 1, "2.3.4.5": 1}
	test_ips := map[string]bool{
		"1.2.3.4": true,
		"2.3.4.5": true,
		"2.3.4.1": false,
		"1.2.4.3": false,
	}

	for ip, r := range test_ips {
		msg := new(dns.Msg)
		s := fmt.Sprintf("example.com. IN A %s", ip)
		rr, err := dns.NewRR(s)
		if err != nil {
			t.Error(err)
		}
		msg.Answer = append(msg.Answer, rr)
		if in_blacklist(msg) != r {
			t.Errorf("%s must match in %v result\n", ip, r)
		}
	}
}
