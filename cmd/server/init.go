package main

import (
	"flag"
	"log"
	"net"
	"strings"
)

var inAddrs []string
var logAddr string
var logConn *net.Conn

func init() {
	in := flag.String("in", "", "Receiving address list. To specify multiple addresses, delimit each port with ' '.")
	lOut := flag.String("log", "", "Specify an address to send log data to.")
	flag.Parse()

	inAddrs = strings.Split(*in, " ")
	logAddr = *lOut

	if len(inAddrs) == 1 && inAddrs[0] == "" {
		log.Println("[Server args]: No inbound addresses declared by user.")
		return
	}
}
