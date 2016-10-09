package main

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"
)

// MyIP my ip
type MyIP struct {
	IP string `json:"origin"`
}

func myip() string {
	res, err := http.Get("https://www.simicloud.com/media/httpbin/ip")
	if err != nil {
		return ""
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	ip := MyIP{}
	err = json.Unmarshal(data, &ip)
	if err != nil {
		return ""
	}
	return ip.IP

}
func TestQuery(t *testing.T) {

	m := net.IPv4Mask(255, 255, 255, 0)
	ip1 := net.ParseIP(myip())
	ipnet := net.IPNet{ip1.Mask(m), m}
	r, err := query("www.taobao.com", "a", ipnet.String(), "", "")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", r)

	r, err = query("www.taobao.com", "a", "", "", "")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", r)

}

func TestHttpsQuery(t *testing.T) {
	m := new(dns.Msg)
	m.SetQuestion("www.taobao.com", dns.TypeA)
	m2, _, err := DefaultHTTPDnsClient.Exchange(m, ServerAddr)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(m2)

	time.Sleep(1 * time.Second)
	m2, _, err = DefaultHTTPDnsClient.Exchange(m, ServerAddr)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(m2)
}

func TestGetMyIP(t *testing.T) {
	a := DefaultHTTPDnsClient.getMyIP()
	time.Sleep(4 * time.Second)
	a = DefaultHTTPDnsClient.getMyIP()
	if a == "" {
		t.Errorf("get ip failed")
	}
}
