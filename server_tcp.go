package main

import (
	"net"

	log "github.com/fangdingjun/go-log/v5"
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

		m, err := getResponseFromUpstream(msg, srv.upstreams)
		if err != nil {
			log.Debugln("query", msg.Question[0].String(), "timeout")
			break
		}

		for _, a := range m.Answer {
			log.Debugln("result", a.String())
		}
		if err := conn.WriteMsg(m); err != nil {
			log.Errorln(err)
			break
		}
	}
}
