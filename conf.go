package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
)

type conf struct {
	Listen          []addr
	BlacklistFile   string
	HostFile        string
	ForwardRules    []rule
	DefaultUpstream []addr
	Timeout         int
	Debug           bool
	blacklist       item
	hosts           hostitem
}

type rule struct {
	Server     []addr
	DomainFile string
	domains    item
}

type addr struct {
	Host    string
	Port    int
	Network string
}

func loadConfig(f string) (*conf, error) {
	c := new(conf)
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}

	if c.Debug {
		logLevel = DEBUG
	}

	if c.blacklist == nil {
		c.blacklist = item{}
	}

	if c.Timeout == 0 {
		c.Timeout = 2
	}

	if err := loadItemFile(c.blacklist, c.BlacklistFile); err != nil {
		return nil, err
	}

	for i := range c.ForwardRules {
		if c.ForwardRules[i].domains == nil {
			c.ForwardRules[i].domains = item{}
		}
		if err := loadItemFile(c.ForwardRules[i].domains,
			c.ForwardRules[i].DomainFile); err != nil {
			return nil, err
		}
	}

	if c.hosts == nil {
		c.hosts = hostitem{}
	}

	if err := loadHostsFile(c.hosts, c.HostFile); err != nil {
		return nil, err
	}

	return c, nil
}

func loadHostsFile(h hostitem, f string) error {
	if f == "" {
		return nil
	}
	fd, err := os.Open(f)
	if err != nil {
		return err
	}
	defer fd.Close()

	r := bufio.NewReader(fd)
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			break
		}
		s1 := strings.Trim(s, " \t\r\n")

		// ignore blank line and comment
		if s1 == "" || s1[0] == '#' {
			continue
		}
		s1 = strings.Replace(s1, "\t", " ", -1)
		s1 = strings.Trim(s1, " \t\r\n")
		ss := strings.Split(s1, " ")

		// ipv4
		t := 1
		if strings.Index(ss[0], ":") != -1 {
			// ipv6
			t = 28
		}

		for _, s2 := range ss[1:] {
			if s2 == "" {
				continue
			}

			h.add(s2, ss[0], t)
		}

	}
	return nil
}

func loadItemFile(it item, f string) error {
	if f == "" {
		return nil
	}
	fd, err := os.Open(f)
	if err != nil {
		return err
	}
	defer fd.Close()

	r := bufio.NewReader(fd)
	for {
		s, err := r.ReadString('\n')
		if s != "" {
			s1 := strings.Trim(s, " \r\n")
			if s1 != "" && s1[0] != '#' {
				it.add(s1)
			}
		}
		if err != nil {
			break
		}
	}
	return nil
}

type item map[string]int

func (it item) has(s string) bool {
	ss := strings.Split(s, ".")

	for i := 0; i < len(ss); i++ {
		s1 := strings.Join(ss[i:], ".")
		if _, ok := it[s1]; ok {
			return true
		}
	}
	return false
}

func (it item) exists(s string) bool {
	_, ok := it[s]
	return ok
}

func (it item) add(s string) {
	it[s] = 1
}

type hostitem map[string][]hostentry

func (ht hostitem) get(domain string, t int) string {
	if v, ok := ht[domain]; ok {
		for _, v1 := range v {
			if v1.domain == domain && v1.t == t {
				return v1.ip
			}
		}
	}
	return ""
}

func (ht hostitem) add(domain, ip string, t int) {
	if v, ok := ht[domain]; ok {
		exists := false
		for _, v1 := range v {
			if v1.domain == domain && v1.ip == ip && v1.t == t {
				exists = true
				break
			}
		}
		if !exists {
			ht[domain] = append(ht[domain], hostentry{domain, ip, t})
		}
	} else {
		v1 := []hostentry{{domain, ip, t}}
		ht[domain] = v1
	}
}

type hostentry struct {
	domain string
	ip     string
	t      int
}
