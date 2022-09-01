package main

import (
	"io"
	"net/http"

	log "github.com/fangdingjun/go-log/v5"
	"github.com/miekg/dns"
)

const dnsMsgType = "application/dns-message"

func (srv *server) handleHTTPReq(w http.ResponseWriter, r *http.Request) {
	/*
		ctype := r.Header.Get("content-type")
		if !strings.HasPrefix(ctype, dnsMsgType) {
			log.Errorf("request type %s, require %s", ctype, dnsMsgType)
			http.Error(w, "dns message is required", http.StatusBadRequest)
			return
		}
	*/

	if r.ContentLength < 10 {
		log.Errorf("message is too small, %v", r.ContentLength)
		http.Error(w, "message is too small", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorln("read request body", err)
		http.Error(w, "read request failed", http.StatusBadRequest)
		return
	}

	msg := new(dns.Msg)
	if err := msg.Unpack(data); err != nil {
		log.Errorln("parse dns message", err)
		http.Error(w, "parse dns message error", http.StatusBadRequest)
		return
	}

	m, err := getResponseFromUpstream(msg, srv.upstreams)
	if err != nil {
		log.Debugln("query", msg.Question[0].String(), "timeout")
		http.Error(w, "query upstream server failed", http.StatusServiceUnavailable)
		return
	}

	for _, a := range m.Answer {
		log.Debugln("result", a.String())
	}
	w.Header().Set("content-type", dnsMsgType)
	d, _ := m.Pack()
	w.Write(d)
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != srv.addr.Path {
		http.Error(w, "Path not found", http.StatusNotFound)
		return
	}
	srv.handleHTTPReq(w, r)
}
