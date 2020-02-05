package main

import (
	"flag"
	"log"
	"strings"
)

var outAddr string
var filenames []string

func init() {
	out := flag.String("out", "", "Specify an address to send to.")
	files := flag.String("files", "", "Outbound files to send. To send multiple files, delimit each filepath with ' '.")
	flag.Parse()

	filenames = strings.Split(*files, " ")
	outAddr = *out

	if outAddr == "" {
		log.Println("[Client Args]: No addresses declared by the user.")
		return
	} else if len(filenames) == 1 && filenames[0] == "" {
		log.Println("[Client Args]: No filenames declared by user.")
		return
	}
}
