package main

import (
	//"fmt"
	"errors"
	"github.com/miekg/dns"
	"net"
	"time"
)

var ip_region []*net.IPNet

type res struct {
	m   *dns.Msg
	err error
}

func _query(m *dns.Msg, s *UpstreamServer, c chan *res) {
	res1 := make(chan *res)

	go query_one(s, m, res1)

	select {
	case r := <-res1:
		c <- r
	case <-time.After(600 * time.Millisecond):
		c <- &res{err: errors.New("timed out")}
	}

}

func query(m *dns.Msg) *dns.Msg {
	resch := make(chan *res, len(DefaultServer))
	for _, s := range DefaultServer {
		go _query(m, s, resch)
	}

	delayed := []*dns.Msg{}
	slen := len(DefaultServer)

	for i := 0; i < slen; i++ {
		r := <-resch
		r1 := *r
		if r1.err != nil {
			logger.Error("error %s\n", r1.err.Error())
			continue
		}

		// drop the result with no error but has an empty result
		if r1.m.Rcode == dns.RcodeSuccess &&
			len(r1.m.Answer) == 0 {
			continue
		}

		// drop blacklist
		if in_blacklist(r1.m) {
			continue
		}

		// check ip region
		if answer_in_region(r1.m, ip_region) {
			return r1.m
		} else {
			delayed = append(delayed, r1.m)
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
