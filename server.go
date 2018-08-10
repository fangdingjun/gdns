package main

import (
	"crypto/tls"
	"net"
	"net/url"
	"strconv"

	"github.com/fangdingjun/go-log"
)

type server struct {
	addr      *url.URL
	cert      string
	key       string
	upstreams []*url.URL
	bootstrap []*url.URL
}

func (srv *server) serve() {
	switch srv.addr.Scheme {
	case "udp":
		srv.serveUDP()
	case "tcp":
		srv.serveTCP()
	case "tls":
		srv.serveTLS()
	case "https":
		srv.serveHTTPS()
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
		log.Fatalf("listen udp error", err)
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
	l, err := tls.Listen("tcp", srv.addr.Host,
		&tls.Config{
			Certificates: []tls.Certificate{cert},
			//NextProtos:   []string{"h2"},
		})
	if err != nil {
		log.Fatalln("listen tls", err)
	}
	defer l.Close()
	log.Debugf("listen tls://%s", l.Addr().String())
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Debugln("tls accept", err)
			break
		}
		go srv.handleTCP(conn)
	}
}

func (srv *server) serveHTTPS() {
	cert, err := tls.LoadX509KeyPair(srv.cert, srv.key)
	if err != nil {
		log.Fatalln("load certificate failed", err)
	}
	l, err := tls.Listen("tcp", srv.addr.Host,
		&tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"h2"},
		})
	if err != nil {
		log.Fatalln("listen https", err)
	}
	defer l.Close()
	log.Debugf("listen https://%s", l.Addr().String())
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Debugln("https accept", err)
			break
		}
		go srv.handleHTTPSConn(conn)
	}
}

func makeServers(c *conf) {
	upstreams := []*url.URL{}
	bootstraps := []*url.URL{}
	for _, a := range c.UpstreamServers {
		u, err := url.Parse(a)
		if err != nil {
			log.Fatal(err)
		}
		upstreams = append(upstreams, u)
	}

	for _, a := range c.BootstrapServers {
		u, err := url.Parse(a)
		if err != nil {
			log.Fatal(err)
		}
		bootstraps = append(bootstraps, u)
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
			bootstrap: bootstraps,
		}
		go srv.serve()
	}
}
