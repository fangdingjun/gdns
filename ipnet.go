package main

import (
	"bufio"
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"os"
	"strings"
	"unicode"
)

// parse ip range region file
// format
//      ip, netmask, prefixlen
func parse_net(fn string) []*net.IPNet {
	fp, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
		return []*net.IPNet{}
	}
	defer fp.Close()

	nets := []*net.IPNet{}

	reader := bufio.NewReader(fp)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.Trim(line, "\r\n")
		fnds := strings.FieldsFunc(line, func(c rune) bool {
			if unicode.IsSpace(c) {
				return true
			}
			if c == ',' {
				return true
			}
			return false
		})

		if len(fnds) != 3 {
			continue
		}

		if _, net1, err := net.ParseCIDR(
			fmt.Sprintf("%s/%s", fnds[0], fnds[2])); err == nil {
			nets = append(nets, net1)

		} else {
			log.Fatal(err)
		}
	}
	return nets
}

// test ip in ip range region
func in_region(ip net.IP, nets []*net.IPNet) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// test dns reply A or AAAA in special ip range region
func answer_in_region(m *dns.Msg, nets []*net.IPNet) bool {
	for _, rr := range m.Answer {
		if a, ok := rr.(*dns.A); ok {
			if in_region(a.A, nets) {
				return true
			}
		}
		if aaaa, ok := rr.(*dns.AAAA); ok {
			if in_region(aaaa.AAAA, nets) {
				return true
			}
		}
	}
	return false
}
