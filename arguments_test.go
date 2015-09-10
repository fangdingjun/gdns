package main

import (
	"testing"
)

func TestArgs(t *testing.T) {
	var a ArgSrvs

	(&a).Set("aa")

	if len(a) != 1 {
		t.Fail()
	}

	if a[0] != "aa" {
		t.Fail()
	}
}
