package main

import (
	"fmt"
	//"github.com/go-yaml/yaml"
	//"io/ioutil"
	"testing"
)

func TestConf(t *testing.T) {
	c, err := loadConfig("testdata/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%#v\n", c)
	if len(c.Listen) != 2 {
		t.Errorf("expected listers 2, got %d", len(c.Listen))
	}

	if !c.blacklist.exists("1.2.3.4") {
		fmt.Printf("%#v\n", c.blacklist)
		t.Errorf("blacklist load failed")
	}
	if c.hosts.get("localhost", 1) != "127.0.0.1" {
		fmt.Printf("%#v\n", c.hosts)
		t.Errorf("hosts load failed")
	}
	if c.hosts.get("localhost", 28) != "::1" {
		fmt.Printf("%#v\n", c.hosts)
		t.Errorf("hosts load failed")
	}
	if len(c.ForwardRules) != 2 {
		fmt.Printf("%#v\n", c.ForwardRules)
		t.Errorf("expected rules 2, got %d", len(c.ForwardRules))
	}
	if !c.ForwardRules[0].domains.has("a.com") {
		fmt.Printf("%#v\n", c.ForwardRules[0].domains)
		t.Errorf("some domains should exit, may be load config failed")
	}
	if !c.ForwardRules[1].domains.has("d.com") {
		fmt.Printf("%#v\n", c.ForwardRules[1].domains)
		t.Errorf("some domains should exit, may be load config failed")
	}
}

func TestItemExists(t *testing.T) {
	it := item{
		"google.cn":     1,
		"www.baidu.com": 1,
		"org":           1,
	}

	testdata := []struct {
		d string
		b bool
	}{
		{"google.cn", true},
		{"www.google.cn", false},
		{"www.a.org", false},
	}

	for _, d := range testdata {
		b1 := it.exists(d.d)
		if b1 != d.b {
			t.Errorf("%s, expected %v, got %v", d.d, d.b, b1)
		}
	}
}

func TestItemHas(t *testing.T) {
	it := item{
		"google.cn":     1,
		"www.baidu.com": 1,
		"org":           1,
	}

	testdata := []struct {
		d string
		b bool
	}{
		{"google.cn", true},
		{"www.google.cn", true},
		{"www.a.org", true},
		{"pan.baidu.com", false},
		{"abc.org", true},
	}

	for _, d := range testdata {
		b1 := it.has(d.d)
		if b1 != d.b {
			t.Errorf("%s, expected %v, got %v", d.d, d.b, b1)
		}
	}
}

func TestItemAdd(t *testing.T) {
	it := item{}
	it.add("www.example.org")
	if !it.exists("www.example.org") {
		t.Errorf("add failed")
	}
}

func TestHostitem(t *testing.T) {
	ht := hostitem{}
	testdata := []hostentry{
		{"www.google.com", "127.0.0.1", 1},
		{"www.google.com", "127.0.0.2", 28},
		{"www.example.org", "127.0.0.3", 28},
		{"www.abc.org", "127.0.0.4", 1},
	}

	for _, d := range testdata {
		ht.add(d.domain, d.ip, d.t)
		ip := ht.get(d.domain, d.t)
		if ip != d.ip {
			t.Errorf("%s, expected %s, got %s", d.domain, d.ip, ip)
		}
	}
	//fmt.Printf("%v\n", ht)
}
