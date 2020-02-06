package main

import (
	"flag"
	"strings"
)

var inAddrs []string
var fileSave bool

func init() {
	in := flag.String("in", "", "Receiving address list. To specify multiple addresses, delimit each port with ' '.")
	save := flag.Bool("save", false, "Saves individual logs for each connection.")
	flag.Parse()

	fileSave = *save
	inAddrs = strings.Split(*in, " ")
}
