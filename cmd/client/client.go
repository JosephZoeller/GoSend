package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

// Transmits files to a target address, one file at a time. Awaits a signal interrupt to terminate.
func main() {
	for _, v := range filenames {
		er := send(v)
		if er != nil {
			log.Printf("Failed to send file %s at [%s] - %s", v, outAddr, er.Error())
		} else {
			log.Printf("Success - File %s sent at [%s].", v, outAddr)
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

// Attempt to speak with an address and, upon connecting, sends the file header information followed by the file.
// Checks to speak with the address each second for 30 seconds, then times out.
func send(filename string) error {
	log.Printf("Client is opening %s", filename)
	fileIn, er := os.Open(filename)
	if er != nil {
		log.Printf("Client failed to open %s - %s", filename, er.Error())
		return er
	}
	log.Printf("Client is attempting to connect at [%s]", outAddr)
	conn, er := connect.SeekConnection(outAddr, 30)
	if er != nil {
		log.Printf("Client failed to connect at [%s]", outAddr)
		return er
	}
	c := *conn
	defer c.Close()
	fHead, er := transit.MakeHeader(fileIn)
	if er != nil {
		log.Println("Client failed to build file header.")
		return er
	}

	er = transit.HeaderOutbound(fHead, conn)
	if er != nil {
		log.Println("Client failed to outbound file header.")
		return er
	}
	log.Printf("Proxy is attempting to outbound file %s.", filename)
	er = transit.FileOutbound(fileIn, conn)
	if er != nil {
		log.Println("Client failed to outbound the file.")
		return er
	}
	return nil
}
