package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type item map[string]int

func (i item) has(s string) bool {
	if _, ok := i[s]; ok {
		return true
	}
	return false
}

func (it item) match(s string) bool {
	iis := strings.Split(s, ".")
	for i := 0; i < len(iis); i++ {
		ii := strings.Join(iis[i:], ".")
		if _, ok := it[ii]; ok {
			return true
		}
	}
	return false
}

type addr struct {
	network string
	addr    string
}

// Rule present a forward rule
type Rule struct {
	DomainlistFile string `json:"domain_list_file"`
	domains        item
	ServersString  []string `json:"servers"`
	servers        []addr
}

type cfg struct {
	Listen         []string `json:"listen"`
	listen         []addr
	ServersString  []string `json:"default_servers"`
	servers        []addr
	Timeout        int      `json:"timeout"`
	BlacklistFiles []string `json:"blacklist_ips"`
	blacklistIps   item
	Rules          []Rule `json:"rules"`
}

func parseCfg(fn string) (*cfg, error) {
	fp, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	c := cfg{}
	buf, err := ioutil.ReadAll(fp)
	err = json.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}

	var adr []addr
	for _, a := range c.ServersString {
		a1 := parseAddr(a)
		if a1.network != "" {
			adr = append(adr, a1)
		}
	}
	c.servers = adr

	var ll []addr
	for _, a := range c.Listen {
		a1 := parseAddr(a)
		if a1.network != "" {
			ll = append(ll, a1)
		}
	}
	c.listen = ll

	l1 := make(item)
	for _, a := range c.BlacklistFiles {
		parseFile(a, &l1)
	}
	c.blacklistIps = l1

	for i, r := range c.Rules {
		l2 := make(item)
		parseFile(r.DomainlistFile, &l2)
		c.Rules[i].domains = l2

		var adr1 []addr
		for _, a := range r.ServersString {
			a1 := parseAddr(a)
			if a1.network != "" {
				adr1 = append(adr1, a1)
			}
		}
		c.Rules[i].servers = adr1
	}
	return &c, nil
}

func parseAddr(addr1 string) addr {
	a := strings.SplitN(addr1, ":", 2)
	if len(a) != 2 {
		fmt.Printf("addr error")
		return addr{"", ""}
	}
	return addr{a[0], a[1]}
}

func parseFile(fn string, i *item) {
	ii := *i
	fp, err := os.Open(fn)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
		return
	}
	defer fp.Close()
	r := bufio.NewReader(fp)
	for {
		line, err := r.ReadString('\n')
		l := strings.Trim(line, " \r\n\t")
		if err != nil {
			if l != "" && l[0] != '#' {
				ii[l] = 1
			}
			break
		}
		if l == "" || l[0] == '#' {
			continue
		}
		ii[l] = 1
	}
}
