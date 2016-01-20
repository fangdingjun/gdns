package main

import (
	//"fmt"
	"github.com/miekg/dns"
	"net"
	"time"
)

var ip_region []*net.IPNet

type res struct {
	m   *dns.Msg
	err error
}

func query(m *dns.Msg) *dns.Msg {
	resch := make(chan *res, len(DefaultServer))
	for _, s := range DefaultServer {
		go query_one(s, m, resch)
	}
	delayed := []*dns.Msg{}
	slen := len(DefaultServer)
	got := 0

loop:
	for {
		select {
		case r := <-resch:
			r1 := *r
			if r1.err != nil {
				logger.Error("error %s\n", r1.err.Error())
				continue
			}
			if in_blacklist(r1.m) {
				continue
			}
			if answer_in_region(r1.m, ip_region) {
				return r1.m
			} else {
				delayed = append(delayed, r1.m)
			}
			got += 1
			if got >= slen {
				break loop
			}
		case <-time.After(900 * time.Millisecond):
			break loop
		}
	}

	if len(delayed) == 0 {
		logger.Error("empty delayed list")
		return nil
	}

	// return first ok result
	for _, m1 := range delayed {
		if m1.Rcode == dns.RcodeSuccess {
			return m1
		}
	}

	// return NXDOMAIN result
	for _, m1 := range delayed {
		if m1.Rcode == dns.RcodeNameError {
			return m1
		}
	}

	// errror
	return nil
}

func query_one(srv *UpstreamServer, m *dns.Msg, ch chan *res) {
	m1, err := srv.query(m)
	select {
	case ch <- &res{m1, err}:
	}
}
