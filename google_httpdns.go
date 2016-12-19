package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ServerAddr is Google dns server ip
var ServerAddr = "74.125.200.100"
var queryIPApi = "https://www.simicloud.com/media/httpbin/ip"

// GoogleHTTPDns struct
type GoogleHTTPDns struct {
	myip string
	l    sync.Mutex
}

func (h *GoogleHTTPDns) getMyIP() string {
	if h.myip != "" {
		return h.myip
	}
	go h.queryMyIP()
	return ""
}

type ipAPI struct {
	IP string `json:"origin"`
}

func (h *GoogleHTTPDns) queryMyIP() {
	h.l.Lock()
	defer h.l.Unlock()
	if h.myip != "" {
		//fmt.Printf("myip: %s\n", h.myip)
		return
	}
	//fmt.Println("get ip...")
	res, err := http.Get(queryIPApi)
	if err != nil {
		//fmt.Println(err)
		return
	}
	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Printf("%s\n", string(d))
	ip := ipAPI{}
	err = json.Unmarshal(d, &ip)
	if err != nil {
		//fmt.Println(err)
		return
	}
	//fmt.Printf("got: %s\n", ip.Ip)
	h.myip = ip.IP
}

func (h *GoogleHTTPDns) getMyNet() string {
	ip := h.getMyIP()
	if ip == "" {
		return ""
	}
	mask := net.IPv4Mask(255, 255, 255, 0)
	ipByte := net.ParseIP(ip)
	ipnet := net.IPNet{ipByte.Mask(mask), mask}
	return ipnet.String()
}

// Exchange send query to server and return the response
func (h *GoogleHTTPDns) Exchange(m *dns.Msg, addr string) (*dns.Msg, time.Duration, error) {
	name := m.Question[0].Name
	t := dns.TypeToString[m.Question[0].Qtype]
	mynet := h.getMyNet()
	r, err := queryGoogleHTTPDNS(name, t, mynet, "", addr)
	if err != nil {
		return nil, 0, err
	}

	m1 := new(dns.Msg)

	m1.SetRcode(m, r.Status)
	for _, rr := range r.Answer {
		_rr := fmt.Sprintf("%s %d IN %s %s", rr.Name, rr.TTL,
			dns.TypeToString[uint16(rr.Type)], rr.Data)

		an, err := dns.NewRR(_rr)
		if err != nil {
			return nil, 0, err
		}
		m1.Answer = append(m1.Answer, an)
	}
	m1.Truncated = r.TC
	m1.RecursionDesired = r.RD
	m1.RecursionAvailable = r.RA
	m1.AuthenticatedData = r.AD
	m1.CheckingDisabled = r.CD
	return m1, 0, nil

}

// Response represent the dns response from server
type Response struct {
	Status           int
	TC               bool
	RD               bool
	RA               bool
	AD               bool
	CD               bool
	Question         []RR
	Answer           []RR
	Additional       []RR
	EDNSClientSubnet string `json:"edns_client_subnet"`
	Comment          string
}

// RR represent the RR record
type RR struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int
	Data string `json:"data"`
}

var httpclient = &http.Client{
	Timeout: 5 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig:     &tls.Config{ServerName: "dns.google.com"},
		TLSHandshakeTimeout: 3 * time.Second,
	},
}

func queryGoogleHTTPDNS(name, t, ednsClientSubnet, padding, srvAddr string) (*Response, error) {
	srvaddr := ServerAddr
	if srvAddr != "" {
		srvaddr = srvAddr
	}
	v := url.Values{}
	v.Add("name", name)
	v.Add("type", t)

	if ednsClientSubnet != "" {
		v.Add("edns_client_subnet", ednsClientSubnet)
	}

	if padding != "" {
		v.Add("random_padding", padding)
	}

	u := fmt.Sprintf("https://%s/resolve?%s", srvaddr, v.Encode())
	r, _ := http.NewRequest("GET", u, nil)
	r.Host = "dns.google.com"
	//r.URL.Host = "dns.google.com"

	res, err := httpclient.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	d := Response{}
	err = json.Unmarshal(data, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}
