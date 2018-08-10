package main

import (
	"net"

	"github.com/fangdingjun/go-log"
	"github.com/miekg/dns"
)

func (srv *server) handleTCP(c net.Conn) {
	defer c.Close()
	log.Debugln("tcp from", c.RemoteAddr())
	conn := dns.Conn{Conn: c}
	for {
		msg, err := conn.ReadMsg()
		if err != nil {
			log.Debugln("tcp read message", err)
			break
		}
		reply := false
		for _, up := range srv.upstreams {
			m, err := queryUpstream(msg, up)
			if err == nil {
				log.Debugln("got reply", m.String())
				conn.WriteMsg(m)
				reply = true
				break
			}
			log.Debugln("query upstream", up.String(), err)
		}
		if !reply {
			break
		}
	}
}
