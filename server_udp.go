package main

import (
	"net"

	log "github.com/fangdingjun/go-log/v5"
	"github.com/miekg/dns"
)

func (srv *server) handleUDP(buf []byte, addr net.Addr, conn *net.UDPConn) {
	msg := new(dns.Msg)
	if err := msg.Unpack(buf); err != nil {
		log.Debugln("udp parse msg", err)
		return
	}

	m, err := getResponseFromUpstream(msg, srv.upstreams)
	if err != nil {
		log.Debugln("query", msg.Question[0].String(), "timeout")
		return
	}

	for _, a := range m.Answer {
		log.Debugln("result", a.String())
	}
	d, _ := m.Pack()
	conn.WriteTo(d, addr)
}
