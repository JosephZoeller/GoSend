package main

import (
	"flag"
	"log"
	"net"
	"strings"
)

var inAddrs []string
var outAddrs []string
var outConns = make([]*net.Conn, 0)
var logAddr string
var logConn *net.Conn
var lastConn *net.Conn

func init() {
	in := flag.String("in", "", "Receiving address list. To specify multiple addresses, delimit each port with ' '.")
	out := flag.String("out", "", "Sending address list. To specify multiple addresses, delimit each port with ' '.")
	lOut := flag.String("log", "", "Specify an address to send log data to.")
	flag.Parse()

	inAddrs = strings.Split(*in, " ")
	outAddrs = strings.Split(*out, " ")
	logAddr = *lOut

	if len(inAddrs) == 1 && inAddrs[0] == "" {
		log.Fatal("[Proxy Args]: No inbound addresses declared by user.")
		return
	} else if len(outAddrs) == 1 && outAddrs[0] == "" {
		log.Fatal("[Proxy Args]: No outbound addresses declared by user.")
		return
	}
}
