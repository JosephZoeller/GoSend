package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
)

var inAddrs []string
var logAddr string
var logConn *net.Conn

func init() {
	var envmsg string

	// log port for the server
	envLogPort := os.Getenv("ServerLogPort")
	if envLogPort != "" {
		envmsg = "Default ServerLogPort = " + envLogPort
	} else {
		envmsg = "There is no default log port specified in the ServerLogPort environment variable."
	}
	logOut := flag.String("log", "", "Specify an address to send log data to. "+envmsg)

	// In ports for server
	envInPorts := os.Getenv("ServerInPorts")
	if envInPorts != "" {
		envmsg = "Default ServerInPorts = " + envInPorts
	} else {
		envmsg = "There are no default inbound ports specified in the ServerInPorts environment variable."
	}
	in := flag.String("in", "", "Receiving address list. To specify multiple addresses, delimit each port with ' '. "+envmsg)

	flag.Parse()

	inAddrs = strings.Split(*in, " ")
	logAddr = *logOut

	if len(inAddrs) == 1 && inAddrs[0] == "" { // splitting "" will result in a slice with 1 ""
		if envInPorts == "" {
			log.Fatal("[Server Args]: No inbound addresses declared by user.")
		}
		inAddrs = strings.Split(envInPorts, " ")
	}
	if logAddr == "" {
		if envLogPort != "" {
			logAddr = envLogPort
		}
	}
}
