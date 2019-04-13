package main

import (
	"bufio"
	"net"

	proxyproto "github.com/pires/go-proxyproto"
)

type protoListener struct {
	net.Listener
}

type protoConn struct {
	net.Conn
	headerDone bool
	r          *bufio.Reader
	proxy      *proxyproto.Header
}

func (l *protoListener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return &protoConn{Conn: c}, err
}

func (c *protoConn) Read(buf []byte) (int, error) {
	var err error
	if !c.headerDone {
		c.r = bufio.NewReader(c.Conn)
		c.proxy, err = proxyproto.Read(c.r)
		if err != nil && err != proxyproto.ErrNoProxyProtocol {
			return 0, err
		}
		c.headerDone = true
		return c.r.Read(buf)
	}
	return c.r.Read(buf)
}

func (c *protoConn) RemoteAddr() net.Addr {
	if c.proxy == nil {
		return c.Conn.RemoteAddr()
	}
	return &net.TCPAddr{
		IP:   c.proxy.SourceAddress,
		Port: int(c.proxy.SourcePort)}
}
