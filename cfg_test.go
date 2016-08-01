package main

import (
	"fmt"
	"os"
	"testing"
)

func TestCfg(t *testing.T) {
	os.Chdir("example_config")
	c, err := parseCfg("config.json")
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	fmt.Printf("%+v\n", c)
	fmt.Printf("%v\n", c.Rules[0].domains.match("google.com"))
	fmt.Printf("%v\n", c.Rules[0].domains.match("www.ip.cn"))
}
