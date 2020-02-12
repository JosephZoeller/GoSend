package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var inAddrs []string
var fileSave bool = true

func init() {
	var envmsg string

	// in ports for the log
	envInPorts := os.Getenv("LogInPorts")
	if envInPorts != "" {
		envmsg = "Default LogInPorts = " + envInPorts
	} else {
		envmsg = "There are no default ports specified in the LogInPorts environment variable."
	}
	in := flag.String("in", "", "Receiving address list. To specify multiple addresses, delimit each port with ' '. "+envmsg)

	// save to log files flag
	//save := flag.Bool("save", false, "Saves individual logs for each connection.") //Not much of a reason to make this a flag.

	flag.Parse()

	//fileSave = *save
	inAddrs = strings.Split(*in, " ")

	if len(inAddrs) == 1 && inAddrs[0] == "" { // splitting "" will result in a slice with 1 ""
		if envInPorts == "" {
			log.Fatal("[Log Manager Args]: No inbound addresses declared by user.")
		}
		inAddrs = strings.Split(envInPorts, " ")
	}
}
