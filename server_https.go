package main

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/fangdingjun/go-log"
	//"github.com/fangdingjun/nghttp2-go"
	"github.com/miekg/dns"
	"golang.org/x/net/http2"
)

func (srv *server) handleHTTPSConn(c net.Conn) {
	defer c.Close()
	tlsconn := c.(*tls.Conn)
	if err := tlsconn.Handshake(); err != nil {
		log.Errorln("handshake", err)
		return
	}
	state := tlsconn.ConnectionState()
	if state.NegotiatedProtocol != "h2" {
		log.Errorln("http2 is needed")
		return
	}
	_srv := &http2.Server{}
	_srv.ServeConn(c, &http2.ServeConnOpts{
		BaseConfig: &http.Server{},
		Handler:    srv,
	})
}

func (srv *server) handleHTTP2Req(w http.ResponseWriter, r *http.Request) {
	ctype := r.Header.Get("content-type")
	if ctype != "application/dns-message" {
		http.Error(w, "dns message is required", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorln("read request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	msg := new(dns.Msg)
	if err := msg.Unpack(data); err != nil {
		log.Errorln("parse dns message", err)
		return
	}
	reply := false
	for _, up := range srv.upstreams {
		log.Debugf("from %s query upstream %s", r.RemoteAddr, up.String())
		log.Debugln("query", msg.Question[0].String())
		m, err := queryUpstream(msg, up)
		if err == nil {
			w.Header().Set("content-type", "application/dns-message")
			w.WriteHeader(http.StatusOK)
			for _, a := range m.Answer {
				log.Debugln("result", a.String())
			}
			d, _ := m.Pack()
			w.Write(d)
			reply = true
			break
		} else {
			log.Errorf("https query upstream %s", err)
		}
	}
	if !reply {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != srv.addr.Path {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	srv.handleHTTP2Req(w, r)
}
