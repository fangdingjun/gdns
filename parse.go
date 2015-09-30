package main

import (
	"encoding/json"
	"errors"
	. "fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"log"
	"net"
	"strconv"
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
	s2 := strings.SplitN(s, ":", 2)

	if len(s2) != 2 {
		msg := Sprintf("error %s not well formatted", s)
		err := errors.New(msg)
		return "", "", err
	}

	if s2[0] != "tcp" && s2[0] != "udp" {
		msg := Sprintf("invalid %s, only tcp or udp allowed", s2[0])
		err := errors.New(msg)
		return "", "", err
	}

	host, port, err := net.SplitHostPort(s2[1])
	if err != nil {
		return "", "", err
	}

	/* check host */
	ip := net.ParseIP(host)
	if ip == nil {
		return "", "", errors.New(Sprintf("invalid host %s", host))
	}

	/* check port */
	_, err = strconv.Atoi(port)
	if err != nil {
		return "", "", err
	}

	return s2[0], s2[1], nil
}

func parse_server(s string) (*UpstreamServer, error) {
	s1 := strings.Split(s, ",")

	if len(s1) != 2 {
		msg := Sprintf("error %s not well formatted", s)
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
