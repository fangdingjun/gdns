package main

import (
	"encoding/json"
	"errors"
	. "fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"log"
	"strings"
)

func load_domain(f string) (Kv, error) {
	var m1 Kv
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(c, &m1)
	if err != nil {
		return nil, err
	}
	return m1, nil
}

func parse_addr(s string) (string, string, error) {
	s2 := strings.Split(s, ":")
	if len(s2) != 3 {
		msg := Sprintf("error %s not well formatted\n", s2)
		err := errors.New(msg)
		return "", "", err
	}
	if s2[0] != "tcp" && s2[0] != "udp" {
		msg := Sprintf("invalid %s, only tcp or udp allowed\n", s2[0])
		err := errors.New(msg)
		return "", "", err
	}
	t := Sprintf("%s:%s", s2[1], s2[2])
	return s2[0], t, nil
}

func parse_server(s string) (*UpstreamServer, error) {
	s1 := strings.Split(s, ",")

	if len(s1) != 2 {
		msg := Sprintf("error %s not well formatted\n", s)
		err := errors.New(msg)
		return nil, err
	}

	proto, addr, err := parse_addr(s1[1])
	if err != nil {
		log.Fatal(err)
	}

	var c *dns.Client
	if proto == "tcp" {
		c = client_tcp
	} else {
		c = client_udp
	}

	d, err := load_domain(s1[0])
	if err != nil {
		log.Print(err)
	}

	var sv *UpstreamServer = &UpstreamServer{
		Addr:    addr,
		domains: d,
		client:  c,
		Proto:   proto,
	}

	return sv, nil
}
