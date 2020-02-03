package connect

import (
	"flag"
	"strings"
)

const inAH string = "Inbound ports to listen on. To open several addresses, delimit each port with ' '."
const outAH string = "Outbound port to send on."
const outFH string = "Outbound files to send. To send multiple files, delimit each filepath with ' '."

func InArgs() []string {
	in := flag.String("in", "", inAH)
	flag.Parse()
	
	return strings.Split(*in, " ")
}

func OutArgs() (string, []string) {
	out := flag.String("out", "", outAH)
	files := flag.String("files", "", outFH)
	flag.Parse()
	
	return *out, strings.Split(*files, " ")
}

func ThroughArgs() ([]string, []string) {
	in := flag.String("in", "", inAH)
	out := flag.String("out", "", outAH)
	flag.Parse()
	
	return strings.Split(*in, " "), strings.Split(*out, " ")
}