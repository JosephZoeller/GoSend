package main

import (
	"errors"
	"log"
	"os"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/logUtil"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

var outAddr string
var filenames []string

func init() {
	outAddr, filenames = connect.OutArgs()
}

// Transmits files to a target address, one file at a time.
func main() {
	if outAddr == "" {
		log.Println(logUtil.FormatError("Cient args", errors.New("No outbound address declared by user.")))
		return
	} else if len(filenames) == 0 {
		log.Println(logUtil.FormatError("Client args", errors.New("No filenames declared by user.")))
		return
	}

	for _, v := range(filenames) {
		er := send(v)
		if er != nil {
			log.Println(er)
			log.Println("Failed to send " + v)
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
		return logUtil.FormatError("Client send", er)
	}
	c := *conn; defer c.Close()

	fHead, er := transit.MakeHeader(fileIn)
	if er != nil {
		return logUtil.FormatError("Client send", er)
	}

	er = transit.HeaderOutbound(fHead, conn)
	if er != nil {
		return logUtil.FormatError("Client send", er)
	}

	er = transit.FileOutbound(fileIn, conn)
	if er != nil {
		return logUtil.FormatError("Client send", er)
	}
	return nil
}