package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/JosephZoeller/gmg/pkg/connect"
	"github.com/JosephZoeller/gmg/pkg/transit"
)

var inAddrs []string
var outAddrs []string
var lastServe = ""

// 1.) Opens the listening addresses (inAddrs)
// 2.) Each listening address awaits a connection
// 3.) Upon connecting with a host, the proxy chooses (round-robin) an address from the speaking addresses (outAddrs)
// 4.) If the speaking address connects to the address, transmit the data. Otherwise, try the next speaking address.
// 5.) After the data is transmitted, the connection is closed and the listener awaits a new connection.

func init() {
	inAddrs, outAddrs = connect.ThroughArgs()
}

// Opens all proxy listeners, then awaits a signal interrupt to terminate.
func main() {
	log.Printf("Registered %d ports", len(inAddrs))
	for i, v := range inAddrs {
		go openListener(v)
		log.Printf("Proxy %d listening at %s", i+1, v)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

// Opens address to listen to and, upon connecting, select a speaking address and reroute the transmission.
func openListener(address string) {

	for {
		lCon, _ := connect.OpenConnection(address)
		for {
			er := sendToAddress(lCon, pickAddress())
			if er == nil {
				break
			}
		}
		c := *lCon
		c.Close()
	}
}

// Reroutes the transmission being sent from the listening connection to the speaking connection.
// Attempts to speak to an address for 5 seconds before getting bored.
func sendToAddress(lCon *net.Conn, sAddr string) error {

	sCon, er := connect.SeekConnection(sAddr, 5)
	if er != nil {
		log.Println(er)
		return er
	}
	log.Println("The reverse proxy has connected to " + sAddr)
	s := *sCon
	defer s.Close()

	fHead, er := transit.PassHeader(lCon, sCon)
	if er != nil {
		log.Println(er)
		return er
	}

	return transit.PassFile(fHead, lCon, sCon)
}

// Load balancer. Selects (via round-robin) an address to speak to from the CLI argument.
func pickAddress() string {
	addrsCnt := len(outAddrs)

	if lastServe == "" {
		lastServe = outAddrs[0]
		log.Println("Choosing " + lastServe + " to serve")
		return lastServe
	}

	for i := 0; i < addrsCnt; i++ {
		if lastServe == outAddrs[i] {
			lastServe = outAddrs[(i+1)%addrsCnt]
			log.Println("Choosing " + lastServe + " to serve")
			break
		}
	}

	return lastServe
}
