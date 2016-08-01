package main

import (
	"fmt"
	"testing"
)

func TestCfg(t *testing.T) {
	c, err := parseCfg("config.json")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%+v\n", c)
	fmt.Printf("%v\n", c.Rules[0].domains.match("google.com"))
	fmt.Printf("%v\n", c.Rules[0].domains.match("www.ip.cn"))
}
