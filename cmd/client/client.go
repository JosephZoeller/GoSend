package main

import (
	"log"
	"os"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

var outAddr string
var filenames []string

func init() {
	outAddr, filenames = connect.OutArgs()
}

// Transmits files to a target address, one file at a time.
func main() {
	for _, v := range(filenames) {
		er := send(v)
		if er != nil {
			log.Println(er)
		}
	}
}

// Attempt to speak with an address and, upon connecting, sends the file header information followed by the file.
// Checks to speak with the address each second for 30 seconds, then times out.
func send(filename string) error {
	fileIn, er := os.Open(filename)
	if er != nil {
		return er
	}

	conn, er := connect.SeekConnection(outAddr, 30)
	if er != nil {
		return er
	}
	c := *conn; defer c.Close()

	fHead, er := transit.MakeHeader(fileIn)
	if er != nil {
		return er
	}

	er = transit.HeaderOutbound(fHead, conn)
	if er != nil {
		return er
	}

	er = transit.FileOutbound(fileIn, conn)
	if er != nil {
		return er
	}
	return nil
}