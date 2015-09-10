package main

import (
	"testing"
)

func TestServerMathNil(t *testing.T) {
	srv := UpstreamServer{} // initial with nil

	domains := []string{"twitter.com", "google.com", "abc.com"}

	for _, d := range domains {
		if srv.match(d) {
			t.Errorf("%s must match in false result\n", d)
		}
	}
}

func TestServerMatch(t *testing.T) {
	d := Kv{"twitter.com": 1, "google.com": 1}
	srv := UpstreamServer{domains: d}

	test_domains := map[string]bool{
		"twitter.com":               true,
		"pbs.twitter.com":           true,
		"abc.pbs.twitter.com":       true,
		"efg.abc.pbs.twitter.com":   true,
		"google.com":                true,
		"plus.google.com":           true,
		"cc.plus.google.com":        true,
		"dd.cc.plus.google.com":     true,
		"twitter.abc.com":           false,
		"twitter.com.aa.com":        false,
		"google.com.cccc.com":       false,
		"google.com.aeddasdfc3.com": false,
	}

	for d, r := range test_domains {
		if srv.match(d) != r {
			t.Errorf("%s must match in %v result\n", d, r)
		}
	}
}
