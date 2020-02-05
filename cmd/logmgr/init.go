package main

import (
	"flag"
	"strings"
)

var inAddrs []string

func init() {
	in := flag.String("in", "", "Receiving address list. To specify multiple addresses, delimit each port with ' '.")
	flag.Parse()

	inAddrs = strings.Split(*in, " ")
}
