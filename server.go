package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/fangdingjun/go-log/v5"
	"github.com/fangdingjun/protolistener"
)

type server struct {
	addr      *url.URL
	cert      string
	key       string
	upstreams []*url.URL
}

func (srv *server) serve() {
	switch srv.addr.Scheme {
	case "udp":
		srv.serveUDP()
	case "tcp":
		srv.serveTCP()
	case "tls":
		srv.serveTLS()
	case "http":
		srv.serveHTTP()
	case "https":
		srv.serveHTTP()
	default:
		log.Fatalf("unsupported type %s", srv.addr.Scheme)
	}
}

func (srv *server) serveUDP() {
	ip, port, _ := net.SplitHostPort(srv.addr.Host)
	_ip := net.ParseIP(ip)
	_port, _ := strconv.Atoi(port)

	udpconn, err := net.ListenUDP("udp", &net.UDPAddr{IP: _ip, Port: _port})
	if err != nil {
		log.Fatalf("listen udp error %s", err)
	}

	defer udpconn.Close()

	buf := make([]byte, 4096)
	for {
		n, addr, err := udpconn.ReadFrom(buf)
		if err != nil {
			log.Debugln(err)
			break
		}
		buf1 := make([]byte, n)
		copy(buf1, buf[:n])
		go srv.handleUDP(buf1, addr, udpconn)
	}
}

func (srv *server) serveTCP() {
	l, err := net.Listen("tcp", srv.addr.Host)
	if err != nil {
		log.Fatalln("listen tcp", err)
	}
	defer l.Close()
	log.Debugf("listen tcp://%s", l.Addr().String())
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Debugln(err)
			break
		}
		go srv.handleTCP(conn)
	}
}

func (srv *server) serveTLS() {
	cert, err := tls.LoadX509KeyPair(srv.cert, srv.key)
	if err != nil {
		log.Fatalln("load certificate failed", err)
	}

	l, err := net.Listen("tcp", srv.addr.Host)
	if err != nil {
		log.Fatalln("listen tls", err)
	}
	defer l.Close()

	log.Debugf("listen tls://%s", l.Addr().String())
	tl := tls.NewListener(protolistener.New(l), &tls.Config{
		Certificates: []tls.Certificate{cert},
		//NextProtos:   []string{"h2"},
	})

	for {
		conn, err := tl.Accept()
		if err != nil {
			log.Debugln("tls accept", err)
			break
		}
		go srv.handleTCP(conn)
	}
}

func (srv *server) serveHTTP() {
	log.Debugf("listen %s://%s", srv.addr.Scheme, srv.addr.Host)

	l, err := net.Listen("tcp", srv.addr.Host)
	if err != nil {
		log.Fatalln("listen https", err)
	}
	defer l.Close()

	httpsrv := &http.Server{
		Handler: LogHandler(srv),
	}
	switch srv.addr.Scheme {
	case "http":
		err = httpsrv.Serve(protolistener.New(l))
	case "https":
		err = httpsrv.ServeTLS(protolistener.New(l), srv.cert, srv.key)
	default:
		err = fmt.Errorf("not supported %s", srv.addr.Scheme)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func makeServers(c *conf) {
	upstreams := []*url.URL{}

	for _, a := range c.UpstreamServers {
		u, err := url.Parse(a)
		if err != nil {
			log.Fatal(err)
		}
		upstreams = append(upstreams, u)
	}

	for _, l := range c.Listen {
		u, err := url.Parse(l.Addr)
		if err != nil {
			log.Fatal(err)
		}
		srv := &server{
			addr:      u,
			cert:      l.Cert,
			key:       l.Key,
			upstreams: upstreams,
		}
		go srv.serve()
	}
}
