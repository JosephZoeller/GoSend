package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var outAddr string
var filenames []string

func init() {
	var envmsg string

	// out port for the client
	envOutPort := os.Getenv("ClientOutPort")
	if envOutPort != "" {
		envmsg = "Default ClientOutPort = " + envOutPort
	} else {
		envmsg = "There is no default port specified in the ClientOutPort environment variable."
	}
	out := flag.String("out", "./test.txt", "Specify an address to send to. "+envmsg)

	// files to send
	files := flag.String("files", "", "Outbound files to send. To send multiple files, delimit each filepath with ' '.")

	
	flag.Parse()

	filenames = strings.Split(*files, " ")
	outAddr = *out

	if outAddr == "" {
		if envOutPort == "" {
			log.Fatal("[Client Args]: No addresses declared by the user.")
		}
		outAddr = envOutPort
	}
	if len(filenames) == 1 && filenames[0] == "" { // splitting "" will result in a slice with 1 ""
		log.Fatal("[Client Args]: No filenames declared by user.")
	}
}

func getDirFiles() { // potentially default transfer files to whatever's in the current directory?
	ioutil.ReadDir(".")
}
