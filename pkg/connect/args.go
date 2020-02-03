package connect

import (
	"flag"
	"strings"
)

const inManyH string = "Receiving address list. To specify multiple addresses, delimit each port with ' '."
const outManyH string = "Sending address list. To specify multiple addresses, delimit each port with ' '."
const outOneH string = "Specify an address to send to."
const outFileH string = "Outbound files to send. To send multiple files, delimit each filepath with ' '."

// CLI argument for specifying a list of addresses to recieve data from. Parses a slice of addresses.
func InArgs() []string {
	in := flag.String("in", "", inManyH)
	flag.Parse()
	
	return strings.Split(*in, " ")
}

// CLI argument for specifying an address to send files to. Parses an address and a slice of filenames.
func OutArgs() (string, []string) {
	out := flag.String("out", "", outOneH)
	files := flag.String("files", "", outFileH)
	flag.Parse()
	
	return *out, strings.Split(*files, " ")
}

// CLI argument for specifying a list of addresses to receive data from, and a list of addresses to send data to. Parses two slices of addresses.
func ThroughArgs() ([]string, []string) {
	in := flag.String("in", "", inManyH)
	out := flag.String("out", "", outManyH)
	flag.Parse()
	
	return strings.Split(*in, " "), strings.Split(*out, " ")
}