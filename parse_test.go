package main

import (
	"testing"
)

func TestParseAddr(t *testing.T) {
	var err error
	var s1, s2 string
	s1, s2, err = parse_addr("udp:123.2.3.4:321")
	t.Logf("parse result %s, %s\n", s1, s2)
	if err != nil {
		t.Fail()
	}

	s1, s2, err = parse_addr("tcp:123.2.3.4:321")
	t.Logf("parse result %s, %s\n", s1, s2)
	if err != nil {
		t.Fail()
	}

	_, _, err = parse_addr("1.2.3.4:333")
	t.Log(err)
	if err == nil {
		t.Fail()
	}

	_, _, err = parse_addr("cc:1.2.3.4:33")
	t.Log(err)
	if err == nil {
		t.Fail()
	}

	_, _, err = parse_addr("tcp:1.2.3.4:33:33")
	t.Log(err)
	if err == nil {
		t.Fail()
	}
}

func TestParseServer(t *testing.T) {
	_, err := parse_server("aa.txt:tcp:1.2.3.4:32")
	t.Log(err)
	if err == nil {
		t.Fail()
	}
	sv, err := parse_server("noexists.txt,tcp:1.2.3.4:32")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if sv.Addr != "1.2.3.4:32" {
		t.Fail()
	}
	if sv.Proto != "tcp" {
		t.Fail()
	}
	if sv.domains != nil {
		t.Fail()
	}
}
