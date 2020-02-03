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

// 1.) open the listeners listed in the json
// 2.) each listener will wait for a connection
// 3.) after a connection, round robin the server ports
// 4.) when a connection ends, the listener is available for another connection

func init() {
	inAddrs, outAddrs = connect.ThroughArgs()
}

func main() {
	log.Printf("Registered %d ports", len(inAddrs))
	for i, v := range inAddrs {
		go OpenListener(v)
		log.Printf("Proxy %d listening at %s", i+1, v)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
}

func OpenListener(address string) {

	ln, _ := net.Listen("tcp", address)
	for {
		lCon, _ := ln.Accept()
		for {
			er := DoConnect(&lCon, pickServe())
			if er == nil {
				break
			}
		}
		lCon.Close()
	}
}

func DoConnect(lCon *net.Conn, sAddr string) error {

	sCon, er := connect.SeekConnection(sAddr, 5)
	log.Println("The reverse proxy has connected to " + sAddr)
	if er != nil {
		log.Println(er)
		return er
	}
	s := *sCon; defer s.Close()

	fHead, er := transit.PassHeader(lCon, sCon)
	if er != nil {
		log.Println(er)
		return er
	}

	return transit.PassFile(fHead, lCon, sCon)
}

func pickServe() string {
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
