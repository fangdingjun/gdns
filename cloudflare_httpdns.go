package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

var httpclientCF = &http.Client{
	Timeout: 8 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig:     &tls.Config{ServerName: "cloudflare-dns.com"},
		TLSHandshakeTimeout: 5 * time.Second,
	},
}

// CloudflareHTTPDns cloudflare http dns
type CloudflareHTTPDns struct {
}

// Exchange send request to server and get result
func (cf *CloudflareHTTPDns) Exchange(m *dns.Msg, addr string) (m1 *dns.Msg, d time.Duration, err error) {
	u := fmt.Sprintf("https://%s/dns-query", addr)
	data, err := m.Pack()
	if err != nil {
		return
	}

	body := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/dns-udpwireformat")

	resp, err := httpclientCF.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	m1 = new(dns.Msg)
	if err = m1.Unpack(data); err != nil {
		return
	}
	return
}
