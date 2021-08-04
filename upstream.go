package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/fangdingjun/go-log/v5"
	"github.com/miekg/dns"
	"golang.org/x/net/http2"
)

var dnsClientTCP *dns.Client
var dnsClientUDP *dns.Client
var dnsClientTLS *dns.Client
var dnsClientHTTPS *httpclient

func getResponseFromUpstream(msg *dns.Msg, upstreams []*url.URL) (*dns.Msg, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resch := make(chan *dns.Msg, len(upstreams))

	for _, up := range upstreams {
		go func(u *url.URL) {
			m, err := queryUpstream(msg, u)
			if err == nil {
				resch <- m
				return
			}
			log.Errorln(u.String(), err)
		}(up)
	}

	var errmsg *dns.Msg

	for i := 0; i < len(upstreams); i++ {
		select {
		case <-ctx.Done():
			return nil, errors.New("time out")
		case m := <-resch:
			if m.MsgHdr.Rcode == dns.RcodeSuccess {
				return m, nil
			}
			errmsg = m
		}
	}
	if errmsg != nil {
		return errmsg, nil
	}
	return nil, errors.New("empty result")
}

func queryUpstream(msg *dns.Msg, upstream *url.URL) (*dns.Msg, error) {
	switch upstream.Scheme {
	case "tcp":
		return queryUpstreamTCP(msg, upstream)
	case "https":
		return queryUpstreamHTTPS(msg, upstream)
	case "udp":
		return queryUpstreamUDP(msg, upstream)
	case "tls":
		return queryUpstreamTLS(msg, upstream)
	default:
	}
	return nil, errors.New("unknown upstream type")
}

func queryUpstreamUDP(msg *dns.Msg, upstream *url.URL) (*dns.Msg, error) {
	m, _, err := dnsClientUDP.Exchange(msg, upstream.Host)
	if err != nil {
		log.Debugf("query udp error %s", err)
	}
	return m, err
}

func queryUpstreamTCP(msg *dns.Msg, upstream *url.URL) (*dns.Msg, error) {
	m, _, err := dnsClientTCP.Exchange(msg, upstream.Host)
	if err != nil {
		log.Debugf("query tcp error %s", err)
	}
	return m, err
}

func queryUpstreamTLS(msg *dns.Msg, upstream *url.URL) (*dns.Msg, error) {
	m, _, err := dnsClientTLS.Exchange(msg, upstream.Host)
	if err != nil {
		log.Debugf("query tls error %s", err)
	}
	return m, err
}

func queryUpstreamHTTPS(msg *dns.Msg, upstream *url.URL) (*dns.Msg, error) {
	m, _, err := dnsClientHTTPS.Exchange(msg, upstream.String())
	if err != nil {
		log.Debugf("query https error %s", err)
	}
	return m, err
}

func initDNSClient(c *conf) {
	var resolver = new(net.Resolver)
	if len(c.BootstrapServers) > 0 {
		log.Debugf("init dns client, bootstrap servers %v", c.BootstrapServers)
		resolver = &net.Resolver{
			Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
				for _, a := range c.BootstrapServers {
					u, _ := url.Parse(a)
					conn, err := net.Dial(u.Scheme, u.Host)
					if err == nil {
						return conn, err
					}
				}
				return nil, errors.New("dial failed")
			},
		}
	}

	dialer := &net.Dialer{
		Resolver: resolver,
	}

	dnsClientTLS = &dns.Client{
		Net:     "tcp-tls",
		Timeout: time.Duration(c.UpstreamTimeout) * time.Second,
		Dialer:  dialer,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: c.UpstreamInsecure,
		},
	}
	dnsClientUDP = &dns.Client{
		Net:     "udp",
		Timeout: time.Duration(c.UpstreamTimeout) * time.Second,
	}
	dnsClientTCP = &dns.Client{
		Net:     "tcp",
		Timeout: time.Duration(c.UpstreamTimeout) * time.Second,
	}
	dnsClientHTTPS = &httpclient{
		Net:     "https",
		Timeout: time.Duration(c.UpstreamTimeout) * time.Second,
		HTTPClient: &http.Client{
			Transport: &http2.Transport{
				DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					log.Debugln("dial to", network, addr)
					p := os.Getenv("https_proxy")
					if p == "" {
						p = os.Getenv("http_proxy")
					}
					if p == "" {
						conn, err := tls.DialWithDialer(dialer, network, addr, cfg)
						return conn, err
					}
					u, _ := url.Parse(p)
					log.Debugf("dial to proxy %s", u.Host)
					conn, err := net.Dial(network, u.Host)
					if err != nil {
						return nil, err
					}
					fmt.Fprintf(conn, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, addr)
					r := bufio.NewReader(conn)
					line, err := r.ReadString('\n')
					if err != nil {
						return nil, err
					}
					ss := strings.Fields(line)
					if ss[1] != "200" {
						return nil, fmt.Errorf("http code %s", ss[1])
					}
					for {
						line, err := r.ReadString('\n')
						if err != nil {
							return nil, err
						}
						line = strings.Trim(line, "\r\n")
						if line == "" {
							break
						}
					}
					tlsconn := tls.Client(conn, cfg)
					if err = tlsconn.Handshake(); err != nil {
						return nil, err
					}
					return tlsconn, nil
				},
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: c.UpstreamInsecure,
					NextProtos:         []string{"h2"},
				},
			},
		},
	}
}
