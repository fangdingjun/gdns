package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/fangdingjun/go-log"
	//nghttp2 "github.com/fangdingjun/nghttp2-go"
	"github.com/miekg/dns"
	"golang.org/x/net/http2"
)

var dnsClientTCP *dns.Client
var dnsClientHTTPS *dns.Client
var dnsClientUDP *dns.Client
var dnsClientTLS *dns.Client

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
	dnsClientHTTPS = &dns.Client{
		Net:     "https",
		Timeout: time.Duration(c.UpstreamTimeout) * time.Second,
		HTTPClient: &http.Client{
			Transport: &http2.Transport{
				DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					log.Debugln("dial to", network, addr)
					conn, err := tls.DialWithDialer(dialer, network, addr, cfg)
					return conn, err
				},
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: c.UpstreamInsecure,
					NextProtos:         []string{"h2"},
				},
			},
		},
	}
}
