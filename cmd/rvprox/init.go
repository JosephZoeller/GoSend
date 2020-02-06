package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
)

var inAddrs []string
var outAddrs []string
var outConns = make([]*net.Conn, 0)
var logAddr string
var logConn *net.Conn
var lastConn *net.Conn

func init() {
	var envmsg string

	// log port for the proxy
	envLogPort := os.Getenv("ProxyLogPort")
	if envLogPort != "" {
		envmsg = "Default ProxyLogPort = " + envLogPort
	} else {
		envmsg = "There is no default log port specified in the ProxyLogPort environment variable."
	}
	logOut := flag.String("log", "", "Specify an address to send log data to. "+envmsg)

	// inbound ports for the proxy
	envInPorts := os.Getenv("ProxyInPorts")
	if envInPorts != "" {
		envmsg = "Default ProxyInPorts = " + envInPorts
	} else {
		envmsg = "There are no default inbound ports specified in the ProxyInPorts environment variable."
	}
	in := flag.String("in", "", "Receiving address list. To specify multiple addresses, delimit each port with ' '. "+envmsg)

	// outbound ports for the proxy
	envOutPorts := os.Getenv("ProxyOutPorts")
	if envOutPorts != "" {
		envmsg = "Default ProxyOutPorts = " + envOutPorts
	} else {
		envmsg = "There are no default outbound ports specified in the ProxyOutPorts environment variable."
	}
	out := flag.String("out", "", "Sending address list. To specify multiple addresses, delimit each port with ' '. "+envmsg)

	
	flag.Parse()

	inAddrs = strings.Split(*in, " ")
	outAddrs = strings.Split(*out, " ")
	logAddr = *logOut

	if len(inAddrs) == 1 && inAddrs[0] == "" { // splitting "" will result in a slice with 1 ""
		if envInPorts == "" {
			log.Fatal("[Proxy Args]: No inbound addresses declared by user.")
		}
		inAddrs = strings.Split(envInPorts, " ")
	}
	if len(outAddrs) == 1 && outAddrs[0] == "" {
		if envOutPorts == "" {
			log.Fatal("[Proxy Args]: No outbound addresses declared by user.")
		}
		outAddrs = strings.Split(envOutPorts, " ")
	}
	if logAddr == "" {
		if envLogPort != "" {
			logAddr = envLogPort
		}
	}
}
