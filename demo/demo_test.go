package main

// Sample failing tests and variables to try.

import (
	"crypto/tls"
	"errors"
	"net"
	"testing"

	"github.com/dbinger/sure"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSame(t *testing.T) {
	be := sure.Be(t)
	be.Same(addr1, addr1)
}

func TestDiff(t *testing.T) {
	be := sure.Be(t)
	be.Diff(nil, 1)
}

var (
	// These variables are here for convenient playground testing.
	err1    = errors.New("basic1")
	err2    = errors.New("basic2")
	errjoin = errors.Join(err1, err2)
	addr1   = net.ParseIP("127.0.0.1")
	addr2   = net.ParseIP("127.0.0.2")
	tls1    = tls.Config{InsecureSkipVerify: true, ServerName: "alpha"}
	tls2    = tls.Config{InsecureSkipVerify: true, ServerName: "beta", MaxVersion: 3}
	tls3    = tls.Config{InsecureSkipVerify: true, ServerName: "gamma"}
	map1    = map[string]int{"A": 1, "B": 2}
	map2    = map[string]int{"A": 1, "B": 2, "C": 3}
	slice1  = []int{1, 2, 3}
	slice2  = []int{3, 2, 1}

	ignoreUnexported = cmpopts.IgnoreUnexported(tls.Config{})
	ignoreServerName = cmpopts.IgnoreFields(tls.Config{}, "ServerName")
	ignoreMaxVersion = cmpopts.IgnoreFields(tls.Config{}, "MaxVersion")
	sortIntSlices    = cmpopts.SortSlices(func(a, b int) bool { return a < b })
	sortStringSlices = cmpopts.SortSlices(func(a, b string) bool { return a < b })
)
