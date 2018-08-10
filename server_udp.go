package main

import (
	"net"

	"github.com/fangdingjun/go-log"
	"github.com/miekg/dns"
)

func (srv *server) handleUDP(buf []byte, addr net.Addr, conn *net.UDPConn) {
	msg := new(dns.Msg)
	if err := msg.Unpack(buf); err != nil {
		log.Debugln("udp parse msg", err)
		return
	}
	for _, up := range srv.upstreams {
		log.Debugf("from %s query upstream %s", addr, up.String())
		log.Debugln("query", msg.Question[0].String())
		m, err := queryUpstream(msg, up)
		if err == nil {
			for _, a := range m.Answer {
				log.Debugln("result", a.String())
			}
			d, _ := m.Pack()
			conn.WriteTo(d, addr)
			break
		} else {
			log.Debugln("udp query upstream err", err)
		}
	}
}
