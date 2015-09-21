package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type HostRecord struct {
	/* RR record */
	rr dns.RR

	/* type, dns.A or dns.AAAA */
	t uint16
}

type Hosts map[string][]HostRecord

/*
   get special type of record form Hosts
*/
func (h Hosts) Get(n string, t uint16) dns.RR {
	n1 := dns.Fqdn(n)
	if hr, ok := h[n1]; ok {
		for _, v := range hr {
			if v.t == t {
				return v.rr
			}
		}
	}
	return nil
}

/*
   read and parse the hosts file
*/
func ReadHosts(fn string) (Hosts, error) {
	fp, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	defer fp.Close()

	hts := Hosts{}

	bf := bufio.NewReader(fp)

	for {
		var t uint16
		bline, _, err := bf.ReadLine()
		if err != nil {
			break
		}

		sline := string(bline)
		sline = strings.Trim(sline, " ")

		/* empty line */
		if sline == "" {
			continue
		}

		/* comment */
		if sline[0] == '#' {
			continue
		}

		/* replace tab to space */
		sline = strings.Replace(sline, "\t", " ", -1)

		lns := strings.Split(sline, " ")

		if len(lns) < 1 {
			return nil, errors.New(fmt.Sprintf("invalid hosts line: %s", sline))
		}

		ip := net.ParseIP(lns[0])
		if ip == nil {
			return nil, errors.New(fmt.Sprintf("invalid ip: %s", lns[0]))
		}

		if strings.Index(lns[0], ".") != -1 {
			t = dns.TypeA
		} else {
			t = dns.TypeAAAA
		}

		for _, dn := range lns[1:] {

			dd := dns.Fqdn(strings.Trim(dn, " "))

			/* ignore space */
			if dd == "." {
				continue
			}

			s := fmt.Sprintf("%s 36000 IN %s %s", dd,
				dns.TypeToString[t], lns[0])

			r, err := dns.NewRR(s)
			if err != nil {
				return nil, err
			}

			if _, ok := hts[dd]; ok {
				hts[dd] = append(hts[dd], HostRecord{r, t})
			} else {
				hts[dd] = []HostRecord{HostRecord{r, t}}
			}
		}
	}

	return hts, nil
}

/*
   return the path of hosts file
*/
func GetHost() string {
	var p string
	if runtime.GOOS == "windows" {
		p = filepath.Join(os.Getenv("SYSTEMROOT"),
			"system32/drivers/etc/hosts")
	} else {
		p = "/etc/hosts"
	}
	return filepath.Clean(p)
}
